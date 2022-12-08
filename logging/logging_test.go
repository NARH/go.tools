package logging

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// ログレベル環境変数名
const ENV_LOG_LEVEL string = "LOG_LEVEL"

// ログ出力するメッセージ
const MESSAGE = "msg"

// 新規ロガー取得試験
func TestNewLogger(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		testLogger := NewLogger()
		if nil == testLogger {
			t.Fail()
		}
		it := reflect.TypeOf((*Logging)(nil)).Elem()
		if !reflect.TypeOf(testLogger).Implements(it) {
			t.Fail()
		}
	})
}

// 環境変数で指定したログレベルの取得試験
func TestGetLogLevel(t *testing.T) {
	/* 試験ケース */
	testCases := []struct {
		/* 試験ケースの構造体 */
		name string // 試験名称
		env  string // 試験コンディション
		want int    // 期待値
	}{
		{
			name: "正常系(未設定)",
			env:  "",
			want: int(INFO),
		},
		{
			name: "正常系(DEBUG)",
			env:  "DEBUG",
			want: int(DEBUG),
		},
		{
			name: "正常系(INFO)",
			env:  "INFO",
			want: int(INFO),
		},
		{
			name: "正常系(WARN)",
			env:  "WARN",
			want: int(WARN),
		},
		{
			name: "正常系(ERROR)",
			env:  "ERROR",
			want: int(ERROR),
		},
		{
			name: "正常系(FATAL)",
			env:  "FATAL",
			want: int(FATAL),
		},
	}

	/* 試験実施 */
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(ENV_LOG_LEVEL, tt.env)
			if result := SetLogLevel(); result != tt.want {
				t.Errorf("SetLogLevel() = %v, want %v", result, tt.want)
			}
		})
	}
}

// enum の文字列表記試験
func TestString(t *testing.T) {
	/* 試験ケース */
	testCases := []struct {
		name string // 試験名称
		enum Level  // 試験対象のenum
		want string // 期待値(enumの文字列表現)
	}{
		{
			name: "正常系(DEBUG)",
			enum: DEBUG,
			want: "DEBUG",
		},
		{
			name: "正常系(INFO)",
			enum: INFO,
			want: "INFO",
		},
		{
			name: "正常系(WARN)",
			enum: WARN,
			want: "WARN",
		},
		{
			name: "正常系(ERROR)",
			enum: ERROR,
			want: "ERROR",
		},
		{
			name: "正常系(FATAL)",
			enum: FATAL,
			want: "FATAL",
		},
		{
			name: "異常系(未定義の値)",
			enum: 99,
			want: "INFO",
		},
	}

	/* 試験実施 */
	var result string

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if result = fmt.Sprint(tt.enum); result != tt.want {
				t.Errorf("String() = %v, want %v", result, tt.want)
			}
		})
	}
}

// FATALレベルのログ試験(Skip)
func TestFatal(t *testing.T) {
	// /* 試験ケース */
	// testCases := []struct {
	// 	name string // 試験名称
	// 	env  string // 試験コンディション
	// 	want bool   // 期待値
	// }{
	// 	{
	// 		name: "正常系(DEBUG)",
	// 		env:  "DEBUG",
	// 		want: true,
	// 	},
	// 	{
	// 		name: "正常系(INFO)",
	// 		env:  "INFO",
	// 		want: true,
	// 	},
	// 	{
	// 		name: "正常系(WARN)",
	// 		env:  "WARN",
	// 		want: true,
	// 	},
	// 	{
	// 		name: "正常系(ERROR)",
	// 		env:  "ERROR",
	// 		want: true,
	// 	},
	// 	{
	// 		name: "正常系(FATAL)",
	// 		env:  "FATAL",
	// 		want: true,
	// 	},
	// }
	//
	// tmpExit := os.Exit
	// defer func() { os.Exit = tmpExit }()
	//
	// for _, tt := range testCases {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		// 環境変数の設定
	// 		t.Setenv(ENV_LOG_LEVEL, tt.env)
	//
	// 		/* 出力先バッファを作成する */
	// 		r, w, _ := os.Pipe()
	//
	// 		logger := NewLogger()
	// 		logger.stdout.SetOutput(w)
	// 		logger.stderr.SetOutput(w)
	//
	// 		logger.Fatal("msg")
	// 		w.Close()
	//
	// 		/* 出力を読み込み */
	// 		var buf bytes.Buffer
	// 		_, err := buf.ReadFrom(r)
	// 		if nil != err {
	// 			t.Fail()
	// 		}
	//
	// 		// 取得値
	// 		msg := strings.TrimRight(buf.String(), "\n")
	//
	// 		// 全ての指定ログレベルに表示される
	// 		if strings.Contains(msg, "FATAL") != tt.want {
	// 			t.Errorf("Fatal() = %v, not contains \"FATAL\" want %v", msg, tt.want)
	// 		}
	// 	})
	// }
	//
}

// ERRORレベルのログ試験
func TestError(t *testing.T) {
	testCases := []logTestCases{
		{
			name: "正常系(DEBUG)",
			env:  "DEBUG",
			want: true,
		},
		{
			name: "正常系(INFO)",
			env:  "INFO",
			want: true,
		},
		{
			name: "正常系(WARN)",
			env:  "WARN",
			want: true,
		},
		{
			name: "正常系(ERROR)",
			env:  "ERROR",
			want: true,
		},
		{
			name: "正常系(FATAL)",
			env:  "FATAL",
			want: false,
		},
	}

	testStdLog("ERROR", testCases, t)
}

// WARNレベルのログ試験
func TestWarn(t *testing.T) {
	testCases := []logTestCases{
		{
			name: "正常系(DEBUG)",
			env:  "DEBUG",
			want: true,
		},
		{
			name: "正常系(INFO)",
			env:  "INFO",
			want: true,
		},
		{
			name: "正常系(WARN)",
			env:  "WARN",
			want: true,
		},
		{
			name: "正常系(ERROR)",
			env:  "ERROR",
			want: false,
		},
		{
			name: "正常系(FATAL)",
			env:  "FATAL",
			want: false,
		},
	}

	testStdLog("WARN", testCases, t)
}

// INFOレベルのログ試験
func TestInfo(t *testing.T) {
	testCases := []logTestCases{
		{
			name: "正常系(DEBUG)",
			env:  "DEBUG",
			want: true,
		},
		{
			name: "正常系(INFO)",
			env:  "INFO",
			want: true,
		},
		{
			name: "正常系(WARN)",
			env:  "WARN",
			want: false,
		},
		{
			name: "正常系(ERROR)",
			env:  "ERROR",
			want: false,
		},
		{
			name: "正常系(FATAL)",
			env:  "FATAL",
			want: false,
		},
	}

	testStdLog("INFO", testCases, t)
}

// DEBUGレベルのログ試験
func TestDebug(t *testing.T) {
	testCases := []logTestCases{
		{
			name: "正常系(DEBUG)",
			env:  "DEBUG",
			want: true,
		},
		{
			name: "正常系(INFO)",
			env:  "INFO",
			want: false,
		},
		{
			name: "正常系(WARN)",
			env:  "WARN",
			want: false,
		},
		{
			name: "正常系(ERROR)",
			env:  "ERROR",
			want: false,
		},
		{
			name: "正常系(FATAL)",
			env:  "FATAL",
			want: false,
		},
	}

	testStdLog("DEBUG", testCases, t)
}

// 出力ログテストケース
type logTestCases struct {
	name string // 試験名称
	env  string // 試験コンディション
	want bool   // 期待値
}

// 各レベルのログ出力を試験する
func testStdLog(level string, testCases []logTestCases, t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			// 環境変数の設定
			t.Setenv(ENV_LOG_LEVEL, tt.env)

			/* 出力先バッファを作成する */
			r, w, _ := os.Pipe()

			// ロガーの取得
			log := NewLogger()
			log.stdout.SetOutput(w)
			log.stderr.SetOutput(w)

			// ログの出力
			switch level {
			case "ERROR":
				log.Error(MESSAGE)
			case "WARN":
				log.Warn(MESSAGE)
			case "INFO":
				log.Info(MESSAGE)
			case "DEBUG":
				log.Debug(MESSAGE)
			}
			w.Close()

			/* 出力を読み込み */
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r)
			if nil != err {
				t.Fail()
			}

			// 取得値
			msg := strings.TrimRight(buf.String(), "\n")

			if strings.Contains(msg, level) != tt.want {
				t.Errorf("Fatal() = %v, not contains \"%s\" want %v", level, msg, tt.want)
			}

		})
	}
}
