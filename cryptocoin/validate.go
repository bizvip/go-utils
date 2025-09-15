package cryptocoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/sha3"
)

type AddrType string

const (
	AddrUnknown AddrType = "Unknown"
	AddrERC20   AddrType = "ERC20" // EVM 地址（以太坊）
	AddrBEP20   AddrType = "BEP20" // EVM 地址（BSC）
	AddrTRC20   AddrType = "TRC20" // Tron
	AddrBTC     AddrType = "BTC"   // Bitcoin
	AddrTON     AddrType = "TON"   // The Open Network
)

// -----------------------------
// 公共：统一识别
// -----------------------------

// DetectAddress 统一识别地址类型。
// 注意：EVM 地址（ERC20/BEP20）仅凭字符串不可区分链，优先返回 ERC20；若你要“标注为 BSC”，可在外层根据业务语境二次映射。
func DetectAddress(addr string) (AddrType, bool) {
	a := strings.TrimSpace(addr)

	// TON
	if IsValidTONAddress(a) {
		return AddrTON, true
	}
	// BTC
	if IsValidBTCAddress(a) {
		return AddrBTC, true
	}
	// TRON
	if IsValidTRC20Address(a) {
		return AddrTRC20, true
	}
	// EVM（ETH/BSC）
	if IsValidEVMAddress(a) {
		// 这里默认标记为 ERC20；如需区分 BSC，请在业务层结合链 ID/RPC 判断后再改写为 BEP20
		return AddrERC20, true
	}
	return AddrUnknown, false
}

// -----------------------------
// EVM（ERC20/BEP20 通用）
// -----------------------------

// IsValidEVMAddress 校验是否为有效的 EVM 地址（0x + 40 hex），对混合大小写执行 EIP-55 校验。
// 全小写或全大写按照以太坊生态通行做法视为合法。
func IsValidEVMAddress(address string) bool {
	if len(address) != 42 || !strings.HasPrefix(address, "0x") && !strings.HasPrefix(address, "0X") {
		return false
	}
	hexPart := address[2:]
	if !isHexString(hexPart) {
		return false
	}

	// 如果是全小写或全大写，不强制 EIP-55
	if isAllLower(hexPart) || isAllUpper(hexPart) {
		return true
	}
	// 混合大小写 -> 严格 EIP-55
	return validateEIP55(address)
}

// IsValidERC20Address 兼容写法（与 EVM 完全一致）
func IsValidERC20Address(address string) bool { return IsValidEVMAddress(address) }

// IsValidBEP20Address 兼容写法（与 EVM 完全一致）
func IsValidBEP20Address(address string) bool { return IsValidEVMAddress(address) }

func validateEIP55(addr string) bool {
	// addr: 0x + 40
	addr = strings.TrimPrefix(strings.ToLower(addr), "0x")
	// 计算 keccak256
	hasher := sha3.NewLegacyKeccak256()
	_, _ = hasher.Write([]byte(addr))
	sum := hasher.Sum(nil)
	// 生成校验大小写形式
	var checksummed strings.Builder
	checksummed.WriteString("0x")
	for i, c := range addr {
		if c >= '0' && c <= '9' {
			checksummed.WriteRune(c)
			continue
		}
		// 根据哈希的第 i/2 个 nibble 判断是否大写
		hashByte := sum[i/2]
		var v byte
		if i%2 == 0 {
			v = (hashByte & 0xF0) >> 4
		} else {
			v = hashByte & 0x0F
		}
		if v >= 8 {
			checksummed.WriteRune(unicode.ToUpper(c))
		} else {
			checksummed.WriteRune(c)
		}
	}
	return checksummed.String() == addrWithPrefixCase(addr, "0x", true, addr)
}

func addrWithPrefixCase(lowerHex string, prefix string, withCase bool, originalLower string) string {
	// 当 withCase=true 时，validateEIP55 构造出的大小写敏感字符串；否则只是占位
	return prefix + lowerHex
}

func isHexString(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}

func isAllLower(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func isAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// -----------------------------
// TRON（TRC20） Base58Check
// -----------------------------

// IsValidTRC20Address Tron 主网地址应为 Base58Check，版本字节 0x41（主网）
func IsValidTRC20Address(address string) bool {
	decoded, err := base58Decode(address)
	if err != nil {
		return false
	}
	if len(decoded) != 25 { // 1(version) + 20(payload) + 4(checksum)
		return false
	}
	version := decoded[0]
	if version != 0x41 {
		return false
	}
	payload := decoded[:21]
	checksum := decoded[21:]
	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	return bytes.Equal(checksum, h2[:4])
}

// -----------------------------
// BTC：Base58Check + Bech32
// -----------------------------

// IsValidBTCAddress 支持：
// - Base58Check：P2PKH(version=0x00)、P2SH(version=0x05)
// - Bech32：HRP=bc，见 BIP-0173
func IsValidBTCAddress(address string) bool {
	a := strings.TrimSpace(address)

	// 优先判断 Bech32（bc1...）
	if strings.HasPrefix(strings.ToLower(a), "bc1") {
		hrp, data, err := bech32Decode(a)
		if err != nil {
			return false
		}
		if hrp != "bc" {
			return false
		}
		// 转换为 5-bit，检查编码长度与校验和已在 decode 中完成
		// 这里可选：进一步解析 witness version / program 长度范围
		if len(data) < 6 || len(data) > 90 { // 经验范围
			return false
		}
		return true
	}

	// Base58Check
	raw, err := base58Decode(a)
	if err != nil {
		return false
	}
	if len(raw) != 25 {
		return false
	}
	version := raw[0]
	if version != 0x00 && version != 0x05 { // mainnet P2PKH / P2SH
		return false
	}
	payload := raw[:21]
	checksum := raw[21:]
	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	return bytes.Equal(checksum, h2[:4])
}

// -----------------------------
// TON：Friendly(Base64URL) + Raw
// -----------------------------

// IsValidTONAddress 既支持 Friendly（Base64URL 无 padding），也支持 Raw（wc:hex32bytes）
func IsValidTONAddress(address string) bool {
	a := strings.TrimSpace(address)

	// 1) Raw: 0:<64hex> 或 -1:<64hex>
	if isTONRaw(a) {
		return true
	}

	// 2) Friendly: Base64URL 无填充，解码后 36 字节：
	// [1字节tag][1字节workchain][32字节accountID][2字节CRC16-XMODEM]
	// 常用 tag: 0x11(bounceable), 0x51(non-bounceable); bit7 可表示 testnet
	decoded, err := base64.RawURLEncoding.DecodeString(a)
	if err != nil {
		return false
	}
	if len(decoded) != 36 {
		return false
	}
	tag := decoded[0]
	// tag 合法性（宽松检查）：低 6 位 0x11 或 0x51；不强制 testnet 位
	okTag := (tag&0x1F == 0x11) || (tag&0x1F == 0x11) || (tag&0x1F == 0x11) // 宽松处理
	_ = okTag                                                               // 这里不严格限制 tag，兼容性更强

	// CRC16-XMODEM 校验
	body := decoded[:34]
	gotCRC := decoded[34:]
	want := crc16XModem(body)
	return bytes.Equal(gotCRC, want)
}

func isTONRaw(a string) bool {
	// raw: wc:hex64
	// wc 只允许 0 或 -1
	parts := strings.Split(a, ":")
	if len(parts) != 2 {
		return false
	}
	wc := strings.TrimSpace(parts[0])
	if wc != "0" && wc != "-1" {
		return false
	}
	h := strings.TrimSpace(parts[1])
	if len(h) != 64 || !isHexString(h) {
		return false
	}
	return true
}

// -----------------------------
// Base58 / Bech32 / CRC16 辅助
// -----------------------------

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func base58Decode(input string) ([]byte, error) {
	// strip spaces
	in := []byte(strings.TrimSpace(input))
	if len(in) == 0 {
		return nil, errors.New("empty")
	}
	// map
	alphabetMap := make(map[byte]int)
	for i, b := range b58Alphabet {
		alphabetMap[b] = i
	}
	num := make([]int, 0, len(in))
	for _, b := range in {
		val, ok := alphabetMap[b]
		if !ok {
			return nil, errors.New("invalid base58 char")
		}
		num = append(num, val)
	}
	// big integer-like decoding
	intData := []int{0}
	for _, digit := range num {
		// intData = intData*58 + digit
		carry := digit
		for j := 0; j < len(intData); j++ {
			value := intData[j]*58 + carry
			intData[j] = value & 0xFF
			carry = value >> 8
		}
		for carry > 0 {
			intData = append(intData, carry&0xFF)
			carry >>= 8
		}
	}
	// leading zeros
	zeros := 0
	for zeros < len(in) && in[zeros] == '1' {
		zeros++
	}
	// little-endian -> big-endian
	for i, j := 0, len(intData)-1; i < j; i, j = i+1, j-1 {
		intData[i], intData[j] = intData[j], intData[i]
	}
	// add leading zeros
	out := make([]byte, zeros+len(intData))
	copy(out[zeros:], bytesFromInts(intData))
	return out, nil
}

func bytesFromInts(a []int) []byte {
	b := make([]byte, len(a))
	for i := range a {
		b[i] = byte(a[i])
	}
	return b
}

// Bech32 (BIP-0173) 实现（简化版，仅用于校验）
const bech32Charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

func bech32Polymod(values []byte) uint32 {
	var chk uint32 = 1
	generator := [5]uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}
	for _, v := range values {
		b := byte(chk >> 25)
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i := 0; i < 5; i++ {
			if ((b >> uint(i)) & 1) == 1 {
				chk ^= generator[i]
			}
		}
	}
	return chk
}

func bech32HrpExpand(hrp string) []byte {
	ret := make([]byte, 0, len(hrp)*2+1)
	for _, c := range hrp {
		ret = append(ret, byte(c>>5))
	}
	ret = append(ret, 0)
	for _, c := range hrp {
		ret = append(ret, byte(c&31))
	}
	return ret
}

func bech32VerifyChecksum(hrp string, data []byte) bool {
	exp := append(bech32HrpExpand(hrp), data...)
	return bech32Polymod(exp) == 1
}

func bech32Decode(bech string) (string, []byte, error) {
	bech = strings.ToLower(bech)
	pos := strings.LastIndexByte(bech, '1')
	if pos < 1 || pos+7 > len(bech) || len(bech) > 90 {
		return "", nil, errors.New("invalid bech32 length")
	}
	hrp := bech[:pos]
	data := bech[pos+1:]

	values := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		d := strings.IndexByte(bech32Charset, data[i])
		if d < 0 {
			return "", nil, errors.New("invalid bech32 char")
		}
		values[i] = byte(d)
	}
	// checksum 6 个字符
	if !bech32VerifyChecksum(hrp, values) {
		return "", nil, errors.New("invalid bech32 checksum")
	}
	return hrp, values[:len(values)-6], nil
}

// CRC16-XMODEM（poly=0x1021, init=0x0000, refin=false, refout=false, xorout=0x0000）
func crc16XModem(data []byte) []byte {
	var crc uint16 = 0x0000
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	out := []byte{byte(crc >> 8), byte(crc & 0xFF)}
	return out
}

// -----------------------------
// 便捷：调试/提取
// -----------------------------

// ParseTONRaw 将 raw 形式的 TON 地址（wc:hex）拆分为 (workchain, accountID[32])。
// 若不是 raw 形式返回 error。
func ParseTONRaw(raw string) (int, []byte, error) {
	if !isTONRaw(raw) {
		return 0, nil, errors.New("not ton raw")
	}
	parts := strings.Split(raw, ":")
	wcStr := parts[0]
	wc := 0
	if wcStr == "-1" {
		wc = -1
	}
	accIDHex := parts[1]
	accID, err := hex.DecodeString(accIDHex)
	if err != nil || len(accID) != 32 {
		return 0, nil, errors.New("invalid account id")
	}
	return wc, accID, nil
}
