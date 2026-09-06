package xjson

import (
	"encoding/json"
	"fmt"
	"github.com/wI2L/jsondiff"
	"testing"
)

func TestDiff(t *testing.T) {
	json1 := `
{
        "created_at": null,
        "description": "我是一个测试虚拟站点",
        "enable": 2,
        "http_type": "https",
        "id": 2,
        "ip": "172.18.100.85",
        "ip_type": 4,
        "mode": 2,
        "name": "测试虚拟站点",
        "port": 8080,
        "proxy_conf": [
            {
                "health_check_id": 1,
                "http_type": "https",
                "id": 1,
                "load_balance_algo": "ip_hash",
                "server_id": 20000,
                "server_name": "server1",
                "servers": [
                    {
                        "ip": "192.168.1.101",
                        "port": 8081,
                        "weight": 1
                    },
                    {
                        "ip": "192.168.1.102",
                        "port": 8081,
                        "weight": 1
                    }
                ],
                "ssl_client_id": 123,
                "ssl_version": [
                    "TLSv1.2",
                    "TLSv1.3"
                ]
            },
            {
                "health_check_id": 2,
                "http_type": "http",
                "id": 2,
                "load_balance_algo": "least_conn",
                "server_id": 20002,
                "server_name": "server2",
                "servers": [
                    {
                        "ip": "192.168.1.103",
                        "port": 8082,
                        "weight": 1
                    },
                    {
                        "ip": "192.168.1.104",
                        "port": 8082,
                        "weight": 1
                    }
                ],
                "ssl_client_id": 456,
                "ssl_version": null
            }
        ],
        "template_acl_id": 0,
        "template_algo_id": 0,
        "template_api_id": 0,
        "template_bot_id": 0,
        "template_dos_id": 0,
        "template_rule_id": 0,
        "updated_at": null
    }
	`
	json2 := `
{
        "created_at": null,
        "description": "我是一个测试虚拟站点",
        "enable": 2,
        "http_type": "https",
        "id": 2,
        "ip": "172.18.100.80",
        "ip_type": 4,
        "mode": 2,
        "name": "测试虚拟站点",
        "port": 8081,
        "proxy_conf": [
            {
                "health_check_id": 1,
                "http_type": "https",
                "id": 1,
                "load_balance_algo": "ip_hash",
                "server_id": 20000,
                "server_name": "server12",
                "servers": [
                    {
                        "ip": "192.168.1.101",
                        "port": 8081,
                        "weight": 1
                    },
                    {
                        "ip": "192.168.1.102",
                        "port": 8081,
                        "weight": 1
                    }
                ],
                "ssl_client_id": 123,
                "ssl_version": [
                    "TLSv1.2",
                    "TLSv1.3"
                ]
            },
            {
                "health_check_id": 2,
                "http_type": "http",
                "id": 2,
                "load_balance_algo": "least_conn",
                "server_id": 20002,
                "server_name": "server2",
                "servers": [
                    {
                        "ip": "192.168.1.103",
                        "port": 8082,
                        "weight": 1
                    },
                    {
                        "ip": "192.168.1.104",
                        "port": 8082,
                        "weight": 1
                    }
                ],
                "ssl_client_id": 456,
                "ssl_version": null
            }
        ],
        "template_acl_id": 0,
        "template_algo_id": 0,
        "template_api_id": 0,
        "template_bot_id": 0,
        "template_dos_id": 0,
        "template_rule_id": 0,
        "updated_at": null
    }
    `
	var obj1, obj2 interface{}
	json.Unmarshal([]byte(json1), &obj1)
	json.Unmarshal([]byte(json2), &obj2)
	patch, err := jsondiff.Compare(obj1, obj2)
	if err != nil {
		// handle error
	}
	print(patch)

	event := map[string]any{
		"old_data": obj1,
		"data":     obj2,
		"diff":     patch,
	}
	print(event)
}

func print(obj any) {
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		// handle error
	}
	fmt.Println(string(b))
}
