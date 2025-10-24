package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

// Scrypt 参数
const (
	N       = 16384 // CPU/Memory cost factor
	R       = 8     // Block size
	P       = 1     // Parallelization factor
	KeyLen  = 32    // Desired key length (in bytes)
	SaltLen = 16    // Length of the salt
)

// GenerateSalt 生成随机盐
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}
	return salt, nil
}

// HashPassword 使用 scrypt 对密码进行加密
func HashPassword(password string) (string, error) {
	// 生成随机盐
	salt, err := GenerateSalt()
	if err != nil {
		return "", err
	}

	// 使用 scrypt 生成密钥
	hashedKey, err := scrypt.Key([]byte(password), salt, N, R, P, KeyLen)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	// 将盐和加密的密钥一起编码为字符串
	saltedHash := append(salt, hashedKey...)
	return base64.StdEncoding.EncodeToString(saltedHash), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解码加密字符串
	data, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %v", err)
	}

	// 提取盐和加密的密钥
	if len(data) < SaltLen+KeyLen {
		return false, fmt.Errorf("invalid hash format")
	}
	salt := data[:SaltLen]
	expectedHash := data[SaltLen:]

	// 使用 scrypt 对输入的密码进行加密
	hashedKey, err := scrypt.Key([]byte(password), salt, N, R, P, KeyLen)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %v", err)
	}

	// 比较两个密钥是否一致（使用 subtle.ConstantTimeCompare 防止时序攻击）
	if subtle.ConstantTimeCompare(hashedKey, expectedHash) == 1 {
		return true, nil
	}
	return false, nil
}

func main() {
	// 要加密的密码
	password := "securepassword123"

	// 1. 加密密码
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Println("Hashed Password:", hashedPassword)

	// 2. 验证密码（正确的密码）
	isValid, err := VerifyPassword(password, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to verify password: %v", err)
	}
	fmt.Println("Password is valid:", isValid)

	// 3. 验证密码（错误的密码）
	wrongPassword := "wrongpassword"
	isValid, err = VerifyPassword(wrongPassword, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to verify password: %v", err)
	}
	fmt.Println("Wrong password is valid:", isValid)
}
