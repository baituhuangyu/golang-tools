package goRSA

import (
    "crypto/md5"
    "fmt"
    "encoding/hex"
    "encoding/base64"
    "strings"
    "time"
    "strconv"
    mrand "math/rand"
)

func scMd5Hex(s string) string {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(s))
    return fmt.Sprintf("%s", hex.EncodeToString(md5Ctx.Sum(nil)))
}

// base64 encode string to sting
func Base64EncodeStr2Str(rawString string) string {

    return base64.StdEncoding.EncodeToString([]byte(rawString))+"\n"
}

// base64 string to sting
func Base64DecodeStr2Str(rawString string) (string, error) {
    debase64result, err := base64.StdEncoding.DecodeString(rawString)
    if err != nil {
        if strings.Contains(rawString, "=="){
            return "", err
        }else {
            return Base64DecodeStr2Str(rawString+"=")
        }
    }
    return string(debase64result), nil
}


func nowUnix() int64 {
    return time.Now().Unix()
}


const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(length int) string {
    mrand.Seed(time.Now().UTC().UnixNano())
    result := make([]byte, length)
    for i := 0; i < length; i++ {
        result[i] = alphanum[mrand.Intn(len(alphanum))]
    }
    return string(result)

    //return alphanum[:length]
}

type AuthCodeSt struct {}

func (a *AuthCodeSt)Random(length int) string {
    return RandomString(length)
}

func (a *AuthCodeSt)QuantumEncode(s string, key string, expiry int64) string {
    return a.AuthCode(s, key, 0, expiry)
}

func (a *AuthCodeSt)QuantumDecode(s string, key string, expiry int64) string {
    return a.AuthCode(s, key, 1, expiry)
}

func (a *AuthCodeSt)AuthCode(s string, key string, operation int, expiry int64) string {
    ENCODE, DECODE := 0, 1

    if len(s) == 0 {
        return ""
    }
    ckey_length := 4
    keya := scMd5Hex(key[:16])
    keyb := scMd5Hex(key[16:][:16])
    keyc := ""
    if ckey_length != 0{
        if operation == DECODE{
            keyc = s[:ckey_length]
        }else if operation == ENCODE{
            keyc = RandomString(ckey_length)
            //keyc = "wa6R"
            fmt.Println("keyc: ", keyc)
        }
    }

    cryptkey := keya + scMd5Hex(keya + keyc)
    key_length := len(cryptkey)
    time_ := ""
    if operation == DECODE{
        s, _ = Base64DecodeStr2Str(s[ckey_length:])
    }else {
        if expiry > 0{
            expiry += nowUnix()
            time_ = fmt.Sprintf("%10d", expiry)
        }else {
            time_ = "0000000000"
        }
        s = time_ + scMd5Hex(s + keyb)[:16] + s
    }

    string_length := len(s)
    result := ""

    rndkey := make([]int, 256)
    for i:=0; i<256; i++  {
        rndkey[i] = int(rune(cryptkey[i % key_length]))
    }

    box := make([]int, 256)
    for i:=0; i<256; i++  {
        box[i] = i
    }

    j := 0
    for i:=0; i<256; i++ {
        j = (j + box[i] + rndkey[i]) % 256
        //box[i], box[j] = box[j], box[i]
        tmp := box[i]
        box[i] = box[j]
        box[j] = tmp
    }
    x, j := 0, 0

    tmpBytes := []byte{}
    for i :=0; i<string_length; i++ {
        x = (x + 1) % 256
        j = (j + box[x]) % 256
        //box[x], box[j] = box[j], box[x]
        tmp := box[x]
        box[x] = box[j]
        box[j]=tmp
        inx := int(rune(s[i])) ^ (box[(box[x] + box[j]) % 256])
        tmpBytes = append(tmpBytes, byte(inx))
    }

    result = string(tmpBytes)

    if operation == DECODE{
        i, _ := strconv.ParseInt(result[:10], 10, 64)
        if (result[:10] == "0000000000" || (i - nowUnix()) > 0) && result[10:26] == scMd5Hex(result[26:] + keyb)[:16]{
            return result[26:]
        }else {
            return ""
        }
    }else {
        return keyc + strings.Replace(Base64EncodeStr2Str(result), "=", "", -1)
    }
}

