package msgc

import "strings"

//# OT 话题判断
func OtMessage(msg string) (bool,string) {
	var ots = [...]string{"java","spring",
		"rust","python","c++","c#","php","notion","android","docker","k8s","kubernetes",
		"javascript","nodejs","node","vue",
		"ssr","v2ray","xray","节点","vpn","科学上网","翻墙","quantumult",
		"qv2ray","shadowsocks","小飞机","qv2rayng","passwall","shadowrocket","小火箭",
		"archlinux","manjaro","ubuntu","kali","渗透","centos"}
	for _, ot := range ots {
		match := strings.Contains(msg, ot)
		if match == true {
			switch ot {
			case "java","spring":
				return true,"@javaer"
			case "rust":
				return true,"@rust_zh"
			case "python":
				return true,"@pythonzh"
			case "c++":
				return true,"@cpluspluszh"
			case "c#":
				return true,"@Csharp_zh"
			case "php":
				return true,"@php_group_cn"
			case "notion":
				return true,"@Notionso"
			case "android":
				return true,"@AndroidDevCn"
			case "docker":
				return true,"@dockertutorial"
			case "k8s","kubernetes":
				return true,"@Kubernetes_CN"
			case "javascript","nodejs","node":
				return true,"@JavaScriptTw"
			case "vue":
				return true,"@vuejs_cn"
			case "ssr","v2ray","xray","节点","vpn","科学上网","翻墙":
				return true,"@v2fly_chat"
			case "qv2ray","shadowsocks","小飞机","qv2rayng","passwall","shadowrocket","小火箭":
				return true,"@Qv2ray_chat"
			case "quantumult":
				return true,"@QuanXApp"
			case "archlinux":
				return true,"@archlinuxcn_group"
			case "manjaro":
				return true,"@manjarolinux_cn"
			case "ubuntu":
				return true,"@ubuntuzh"
			case "kali","渗透":
				return true,"@hackerzh"
			case "centos":
				return true,"@centoszh"
			default:
				return true,"@teapers"
			}

		}
	}
	return false,"@teapers"
}
