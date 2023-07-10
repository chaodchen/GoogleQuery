package api


type UIParameter struct {
	Word  string
	Time  string
	Web   string
	Type   string
	Proxy string
}

type DefaultParameter struct {
    hl              string      // 搜索结果语言偏好
    as_q            string      // 搜索指定关键词
    as_epq          string      // 精确匹配短语
    as_oq           string      // 指定必要关键词
    as_eq           string      // 排除指定关键词
    as_nlo          string      // 数字范围从xx
    as_nhi          string      // 到xx
    lr              string      // 搜索结果的语言限定
    cr              string      // 搜索结果的地区限定
    as_qdr          string      // 搜索结果的时间限定
    as_sitesearch   string      // 指定搜索结果在特定网站内
    as_occt         string      // 指定关键词存在的位置
    safe            string      // 搜索结果安全过滤级别
    as_filetype     string      // 搜索结果的文件类型限定
    tbs             string      // 包含其他高级搜索选项
}

type ProxyList struct {
	NoProxy		string
	Local 		string
	Alpha 		string
	Beta 		string
	Gamma 		string
    
}

type Config struct {

}