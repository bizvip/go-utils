// Package bizid 业务 ID 生成器 + 项目级 snowflake 单例入口。
//
// 业务 ID 格式：<prefix> + base26(YYMMDDHHmm) + <snowflakeID>
//
//   - <prefix>      业务前缀（2 个大写字母，如 RG=充值、WD=提款、DP=存款 ...）
//   - base26        把当前 UTC 时间的「年月日时分」拼成 10 位十进制数字后用 A-Z 编码，
//     得到长度 ≤ 8 的字母串。该段同一分钟内所有 ID 一致，便于"看一眼时间"。
//   - <snowflakeID> 由本包持有的 ShortIdGenerator 单例生成，与裸 snowflake 取自同一实例。
//
// 使用流程：
//
//	import "github.com/bizvip/go-utils/base/id/bizid"
//
//	// 程序启动时调用一次（每个 app 分配独立 workerId，范围 0-31）
//	if err := bizid.Init(5); err != nil { ... }
//
//	// 业务 ID
//	orderNo := bizid.GetRechargeOrderID()  // 例：RGFQRJVPK543304808
//
//	// 原始雪花 ID（int64，可直接做 PG 主键 / 业务整数 ID）
//	id := bizid.MustGetSnowflakeID()
package bizid

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bizvip/go-utils/base/snowflake"
)

var (
	gen *snowflake.ShortIdGenerator
	mu  sync.RWMutex
)

const base26Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ErrNotInitialized 在 Init 调用前使用本包任意 ID 接口时返回 / panic 携带的根错误。
var ErrNotInitialized = errors.New("bizid: Init(workerId) must be called before use")

// Init 创建并保存进程级 ShortIdGenerator 单例。每个 app（adminapi / clientapi / ...）
// 必须在启动入口调用一次，传入该 app 独占的 workerId（0-31，由部署规约分配）。
//
// 重复调用会替换旧实例。已经基于旧实例发出的 ID 不受影响，新调用走新实例。
// **不要**在同一进程外再独立创建 ShortIdGenerator——否则不同实例同 workerId 同毫秒
// 会撞 ID（每个实例只在自身 mutex 内保证序列单调）。
func Init(workerId int64) error {
	g, err := snowflake.NewShortIdGenerator(workerId)
	if err != nil {
		return fmt.Errorf("bizid: init snowflake: %w", err)
	}
	mu.Lock()
	gen = g
	mu.Unlock()
	return nil
}

func currentGen() *snowflake.ShortIdGenerator {
	mu.RLock()
	defer mu.RUnlock()
	return gen
}

// GetSnowflakeID 返回 int64 雪花 ID。本项目 snowflake 总位数 48（time 39 + worker 5 + seq 4），
// 远小于 int63，转换 int64 安全。
func GetSnowflakeID() (int64, error) {
	g := currentGen()
	if g == nil {
		return 0, ErrNotInitialized
	}
	id, err := g.NextID()
	if err != nil {
		return 0, fmt.Errorf("bizid: snowflake next id: %w", err)
	}
	return int64(id), nil
}

// MustGetSnowflakeID 同 GetSnowflakeID，失败时 panic。适合放进 ent DefaultFunc 这种"不能返回 error"的位置。
func MustGetSnowflakeID() int64 {
	id, err := GetSnowflakeID()
	if err != nil {
		panic(err)
	}
	return id
}

// New 生成业务 ID：<prefix> + base26(YYMMDDHHmm) + <snowflakeID>。
// 推荐 prefix 用 1-3 个大写字母。
func New(prefix string) string {
	id := MustGetSnowflakeID()
	return prefix + timeBase26(time.Now().UTC()) + strconv.FormatInt(id, 10)
}

// timeBase26 把 UTC 时间的「年(2)月(2)日(2)时(2)分(2) = 10 位十进制数字」做 base26 编码。
// 输出长度 ≤ 8 字符。
func timeBase26(t time.Time) string {
	yy := uint64(t.Year() % 100)
	digits := yy*100000000 +
		uint64(t.Month())*1000000 +
		uint64(t.Day())*10000 +
		uint64(t.Hour())*100 +
		uint64(t.Minute())
	return encodeBase26(digits)
}

// encodeBase26 把无符号整数编码为 A-Z 字母串（A=0, Z=25）。0 -> "A"。
func encodeBase26(n uint64) string {
	if n == 0 {
		return "A"
	}
	// 最大输入 99_12_31_23_59 = 9_912_312_359，log26 后 ≈ 7.05 → 8 位足够。
	buf := make([]byte, 0, 8)
	for n > 0 {
		buf = append(buf, base26Alphabet[n%26])
		n /= 26
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// ----- 业务快捷函数（前缀约定） -----

// GetRechargeOrderID 充值订单 ID（RG = recharge）
func GetRechargeOrderID() string { return New("RG") }

// GetWithdrawID 提款单 ID（WD = withdraw）
func GetWithdrawID() string { return New("WD") }

// GetDepositOrderID 存款订单 ID（DP = deposit）
func GetDepositOrderID() string { return New("DP") }

// GetBuyOrderID 购买 / 消费订单 ID（BU = buy）
func GetBuyOrderID() string { return New("BU") }

// GetExchangeID 兑换订单 ID（EX = exchange）
func GetExchangeID() string { return New("EX") }

// GetBonusOrderID 积分订单 ID（BO = bonus order）
func GetBonusOrderID() string { return New("BO") }

// GetGemOrderID 宝石订单 ID（GM = gem）
func GetGemOrderID() string { return New("GM") }

// GetGoldOrderID 金币订单 ID（GD = gold）
func GetGoldOrderID() string { return New("GD") }
