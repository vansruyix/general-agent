package errorx

const (
	// DefaultSuccessCode 默认成功响应码
	DefaultSuccessCode = 0
	// DefaultErrorCode 默认错误码（如果没有具体错误码code定义，则使用该值）
	DefaultErrorCode = -1
)

// <错误定义规范>
// 通用错误：与http状态码保持一致
// 客户端错误：40开头
// 登录类错误：401开头
// 服务端错误：50开头

// 默认错误：与http状态码保持一致
var (
	// ErrDefault 未定义错误，这类错误通常需要业务代码自己填充，message内容
	ErrDefault        = NewErrorWithLevel(DefaultErrorCode, "未定义错误", LevelWarning)
	ErrUnknown        = NewErrorWithLevel(500, "服务端内部错误", LevelError)
	ErrRequestParam   = NewErrorWithLevel(400, "请求参数错误", LevelInfo)
	ErrRequestEnParam = NewErrorWithLevel(400, "Request parameter error", LevelInfo)
	ErrForbidden      = NewErrorWithLevel(403, "请求被拒绝", LevelInfo)
	ErrUnauthorized   = NewErrorWithLevel(401, "未登录", LevelInfo)
)

// 登录类错误：401开头
var (
	ErrInvalidAuthorizationCode    = NewErrorWithLevel(401000, "授权码已失效", LevelInfo)
	ErrNotSupportedAuthorization   = NewErrorWithLevel(401001, "未支持的认证类型", LevelInfo)
	ErrUserDisable                 = NewErrorWithLevel(401002, "您的账户已被禁用", LevelInfo)
	ErrUserLocked                  = NewErrorWithLevel(401003, "账户被锁定请稍后再试", LevelInfo)
	ErrLogin                       = NewErrorWithLevel(401004, "登录失败", LevelInfo)
	ErrIncorrectUsernameOrPassword = NewErrorWithLevel(401005, "用户名或密码错误", LevelInfo)
	ErrChangePassword              = NewErrorWithLevel(401006, "需要修改密码", LevelInfo)
	ErrPasswordSimple              = NewErrorWithLevel(401007, "密码太简单", LevelInfo)
	ErrInvalidCaptcha              = NewErrorWithLevel(401008, "验证码错误", LevelInfo)
	ErrInvalidToken                = NewErrorWithLevel(401009, "令牌错误", LevelInfo)
	Err2FA                         = NewErrorWithLevel(401010, "需要进一步认证", LevelInfo)
	ErrExpiredToken                = NewErrorWithLevel(401011, "令牌已过期", LevelInfo)
)

// 通用参数类错误 402 开头

var (
	ErrDuplication     = NewErrorWithLevel(402000, "数据重复", LevelInfo)
	ErrNameDuplication = NewErrorWithLevel(402001, "名称重复", LevelInfo)
)

// 通用业务操作类错误，500 开头
var (
	ErrCreate     = NewErrorWithLevel(500001, "创建失败", LevelInfo)
	ErrUpdate     = NewErrorWithLevel(500002, "更新失败", LevelInfo)
	ErrDelete     = NewErrorWithLevel(500003, "删除失败", LevelInfo)
	ErrFind       = NewErrorWithLevel(500004, "查找失败", LevelInfo)
	ErrNotExists  = NewErrorWithLevel(500005, "记录不存在或已删除", LevelInfo)
	ErrExists     = NewErrorWithLevel(500006, "记录已存在", LevelInfo)
	ErrImport     = NewErrorWithLevel(500007, "导入失败", LevelInfo)
	ErrExport     = NewErrorWithLevel(500008, "导出失败", LevelInfo)
	ErrUsing      = NewErrorWithLevel(500009, "资源在使用中", LevelInfo)
	ErrApply      = NewErrorWithLevel(500010, "应用配置失败", LevelInfo)
	ErrApplyBiz   = NewErrorWithLevel(500011, "应用配置中", LevelInfo)
	ErrInternal   = NewErrorWithLevel(500012, "快照创建失败", LevelInfo)
	ErrInternalHF = NewErrorWithLevel(500012, "快照恢复失败", LevelInfo)
)

// 文件类错误，501 开头
var (
	ErrFileNotExists = NewErrorWithLevel(501001, "文件不存在", LevelInfo)
	ErrReadFile      = NewErrorWithLevel(501002, "文件读取错误", LevelInfo)
	ErrWriteFile     = NewErrorWithLevel(501003, "文件写入错误", LevelInfo)
	ErrParseFile     = NewErrorWithLevel(501004, "解析文件失败", LevelInfo)
	ErrDownload      = NewErrorWithLevel(501005, "下载失败", LevelInfo)
)
var (
	ErrUpgrade = NewErrorWithLevel(502001, "升级失败", LevelInfo)
)

// 远程调用错误，503开头
var (
	ErrDeviceNotExists = NewErrorWithLevel(503001, "设备不存在或已离线", LevelInfo)
	ErrDeviceRpcError  = NewErrorWithLevel(503002, "设备远程调用失败", LevelInfo)
	ErrDeviceIdInvalid = NewErrorWithLevel(503003, "设备ID错误", LevelInfo)
)

// ErrCert 证书类错误，504开头
var (
	ErrCert = NewErrorWithLevel(504001, "证书错误", LevelInfo)
)

var (
	ErrRpcError = NewErrorWithLevel(505001, "远程调用失败", LevelError)
)

var (
	ErrProviderNotFound   = NewErrorWithLevel(506001, "模型供应商不存在", LevelError)
	ErrProviderNotSupport = NewErrorWithLevel(506002, "暂不支持该模型供应商", LevelError)
)
