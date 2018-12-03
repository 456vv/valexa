package valexa

import (
	"net/http"
	"io"
	"fmt"
	"net/url"
	"path"
	"strings"
	"bytes"
	"io/ioutil"
	"os"
	"encoding/pem"
	"encoding/base64"
	"encoding/json"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/rsa"
	"time"

)



//https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#hosting-a-custom-skill-as-a-web-service
type checkRequestBody struct {
	R 				*http.Request							//请求对象
}

func (T *checkRequestBody) echoRequest(echoApp *EchoApplication) (echoReq *EchoRequest, err error) {
	err = json.NewDecoder(T.R.Body).Decode(&echoReq)
	if err != nil {
		return nil, fmt.Errorf("valexa: Body 数据结构无法解析， 错误的是（%s）", err)
	}

	//检查时间
	if !echoReq.VerifyTimestamp(echoApp.ValidReqTimestamp) {
		return nil, fmt.Errorf("valexa: 请求时间超出（>%ds），已经过时了。", echoApp.ValidReqTimestamp)
	}
	return echoReq, nil

}

func (T *checkRequestBody) verifyBody(echoApp *EchoApplication) (body io.Reader, err error) {

	if T.R.Method != "POST" {
		return nil, fmt.Errorf("valexa: 请求仅支持 POST 方法， 错误的是（%s）", T.R.Method)
	}

	certURL := T.R.Header.Get("SignatureCertChainUrl")

	link, err := url.Parse(certURL)
	if err != nil{
		return nil, fmt.Errorf("valexa: 解析SignatureCertChainUrl地址路径失败， 错误的是（%s）", err)
	}

	if !strings.EqualFold(link.Scheme, "https") {
		return nil, fmt.Errorf("valexa: 网址协议仅支持https， 错误的是（%s）", link.Scheme)
	}

	if !strings.EqualFold(link.Host, "s3.amazonaws.com")  && !strings.EqualFold(link.Host, "s3.amazonaws.com:443") {
		return nil, fmt.Errorf("valexa: 网址host仅支持s3.amazonaws.com， 错误的是（%s）", link.Host)
	}

	if !strings.HasPrefix(path.Clean(link.Path) , "/echo.api/") {
		return nil, fmt.Errorf("valexa: 网址Path前缀仅支持/echo.api/， 错误的是（%s）", link.Path)
	}

	//读取证书文件
	name 		:= path.Base(link.Path)
	filePath 	:= path.Join(echoApp.CertFolder, name)
	certBody, err := ioutil.ReadFile(filePath)
	if err != nil {
		resp, err := http.Get(certURL)
		if err != nil {
			return nil, fmt.Errorf("valexa: 下载证书文件失败， 错误的是（%s）", err)
		}
		certBody, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("valexa: 读取文件失败， 错误的是（%s）", err)
		}
		osFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, fmt.Errorf("valexa: 创建文件失败， 错误的是（%s）", err)
		}
		n, err := osFile.Write(certBody)
		osFile.Close()
		if err != nil {
			return nil, fmt.Errorf("valexa: 写入文件失败， 错误的是（%s）", err)
		}
		if len(certBody) != n {
			os.Rename(filePath, filePath+".temp")
			return nil, fmt.Errorf("valexa: 证书文件保存到本地不完整！")
		}
	}
	if len(certBody) == 0 {
		return nil, fmt.Errorf("valexa: 证书文件大小为 0")
	}

	//如果不能识别这个证书，需要重命名
	var rename bool = true
	defer func(){
		if err != nil && rename {
			os.Rename(filePath, filePath+".temp")
		}
	}()

	var (
		cCert	*x509.Certificate
		rCert	*x509.Certificate
	)

	pemBlock, certBody := pem.Decode(certBody)
	if pemBlock == nil {
		return nil, fmt.Errorf("valexa: 无法解析证书PEM文件！")
	}

	x509Certificate, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("valexa: 无法解析证书PEM， 错误的是（%s）", err)
	}
	cCert = x509Certificate
	for len(certBody)>0 {
		pemBlock, certBody = pem.Decode(certBody)
		if pemBlock == nil {
			return nil, fmt.Errorf("valexa: 无法解析证书PEM文件！")
		}
		rCert, err = x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("valexa: 无法解析证书PEM， 错误的是（%s）", err)
		}
		if err := cCert.CheckSignatureFrom(rCert); err != nil {
			return nil, fmt.Errorf("valexa: 无法验证证书链签名， 错误的是（%s）", err)
		}
		cCert = rCert
	}
//	Amazon 提供的证书链是不完整的，无法使用根证书验证自身
//	所以这里注释
//	if err := cCert.CheckSignatureFrom(cCert); err != nil {
//		return nil, fmt.Errorf("根证书无法验证自身签名， 错误的是（%s）", err)
//	}


	if time.Now().Unix() < x509Certificate.NotBefore.Unix() || time.Now().Unix() > x509Certificate.NotAfter.Unix() {
		return nil, fmt.Errorf("valexa: Amazon 证书已经过期！")
	}

	//检查证书签属名称
	foundName := false
	for _, altName := range x509Certificate.Subject.Names {
		if altName.Value.(string) == "echo-api.amazon.com" {
			foundName = true
			break
		}
	}
	if !foundName {
		return nil, fmt.Errorf("valexa: Amazon 证书 Subject.names[].Value 没有检测到包含 echo-api.amazon.com 域名。")
	}

	//如果错误，不要命名证书文件
	rename = false

	//验证KEY
	publicKey := x509Certificate.PublicKey
	encryptedSig, err := base64.StdEncoding.DecodeString(T.R.Header.Get("Signature"))
	if err != nil {
		return nil, fmt.Errorf("valexa: 请求标头 Signature 无法识别， 错误的是（%s）", T.R.Header.Get("Signature"))
	}

	//读取Body, 和转化 HASH
	var bodyBuf bytes.Buffer
	hash := sha1.New()
	ioReader := io.TeeReader(T.R.Body, &bodyBuf)
	_, err = io.Copy(hash, ioReader)
	T.R.Body.Close()
	T.R.Body = ioutil.NopCloser(&bodyBuf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("valexa: 读取 Body 数据转化成 sha1 HASH 出了问题， 错误的是（%s）", err)
	}
	
	if err := rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig); err != nil {
		return nil, fmt.Errorf("valexa: 证书无法验证 Body 数据， 错误的是（%s）", err)
	}
	
	
	return &bodyBuf, nil
}




















