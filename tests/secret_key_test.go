package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func checkValidDevice(secret string, t *testing.T) (string, bool) {
	privateKeyPath := "/home/zht/.ssh/id_rsa"
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Errorf("Failed to read private key file: %v", err)
		return "", false
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		t.Errorf("Failed to parse private key: %v", err)
		return "", false
	}

	// Base64解码secret参数
	encryptedData, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		t.Errorf("Failed to decode base64 secret: %v", err)
		return "", false
	}

	// 用私钥解密
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), encryptedData)
	if err != nil {
		t.Errorf("Failed to decrypt data: %v", err)
		return "", false
	}

	secretList := strings.Split(string(decryptedData), "|")
	if secretList[0] != "test" {
		t.Logf("Invalid secret prefix: %s", secretList[0])
		return "", false
	}

	return secretList[1], true
}

func TestSecretKey(t *testing.T) {
	secret := "test|FF-FF-FF-FF-FF-FF"

	// 读取公钥文件
	publicKeyPath := "/home/zht/.ssh/id_rsa.pub"
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to read public key file: %v", err)
	}

	// 解析SSH公钥格式
	publicKeyStr := strings.TrimSpace(string(publicKeyBytes))
	parts := strings.Fields(publicKeyStr)
	if len(parts) < 2 {
		t.Fatal("Invalid SSH public key format")
	}

	// 解码base64编码的公钥数据
	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// 解析SSH公钥
	publicKey, err := x509.ParsePKCS1PublicKey(keyData[len("ssh-rsa")+4:])
	if err != nil {
		// 如果PKCS1失败，尝试使用ssh包解析
		t.Logf("PKCS1 parse failed, trying alternative method: %v", err)

		// 使用golang.org/x/crypto/ssh包解析SSH公钥
		sshPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
		if err != nil {
			t.Fatalf("Failed to parse SSH public key: %v", err)
		}

		cryptoPublicKey := sshPublicKey.(ssh.CryptoPublicKey).CryptoPublicKey()
		var ok bool
		publicKey, ok = cryptoPublicKey.(*rsa.PublicKey)
		if !ok {
			t.Fatal("Public key is not RSA")
		}
	}

	// 使用公钥加密secret
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(secret))
	if err != nil {
		t.Fatalf("Failed to encrypt secret: %v", err)
	}

	// Base64编码加密后的数据
	encryptedSecret := base64.StdEncoding.EncodeToString(encryptedData)
	t.Logf("Encrypted secret: %s", encryptedSecret)

	// 测试解密
	decryptedSecret, valid := checkValidDevice(encryptedSecret, t)
	if !valid {
		t.Fatal("Failed to validate encrypted secret")
	}

	t.Logf("Decrypted MAC: %s", decryptedSecret)

	if decryptedSecret != "FF-FF-FF-FF-FF-FF" {
		t.Fatalf("Expected MAC 'FF-FF-FF-FF-FF-FF', got '%s'", decryptedSecret)
	}

	t.Log("Secret key encryption/decryption test passed")
}
