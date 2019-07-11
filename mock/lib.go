/*
Package mock will mock tcp/udp endpoints
*/
package mock

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"sync"

	"github.com/psyomn/psy/common"

	"github.com/go-yaml/yaml"
)

const example = `---
some-service:
  type: tcp
  port: 9999
  return: nil # for no returns

some-service-2:
  type: tcp
  port: 9998
  return: "the text"

byte-tcp-service-3:
  type: tcp
  port: 9997
  return: 13,14,15

udp-service:
  type: udp
  port: 9996
  return: "blah"

udp-service-bytes:
  type: udp
  port: 9995
  return: 12,13,14
`

type record struct {
	Type   string      `yaml:"type"`
	Port   int         `yaml:"port"`
	Return interface{} `yaml:"return"`
}

type config map[string]record

func printUsage() {
	fmt.Println("usage: ")
	fmt.Println("  mock [--generate] config.yaml")
	fmt.Println("       --generate will generate a sample config")
}

// Run net mocker
func Run(args common.RunParams) common.RunReturn {
	t := &config{}

	if len(args) == 0 {
		printUsage()
		return errors.New("wrong usage of mock")
	}

	if len(args) >= 2 && args[0] == "--generate" {
		return generateYamlConfig(args[1])
	}

	if len(args) != 1 {
		printUsage()
		return errors.New("wrong usage of mock")
	}

	configContents, err := readFile(args[0])
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configContents, t)
	if err != nil {
		return err
	}

	return processEntries(t)
}

func processEntries(conf *config) error {
	var wg sync.WaitGroup

	for _, v := range *conf {
		switch v.Type {
		case "udp":
			wg.Add(1)
			go createUDP(v.Port, v.Return, &wg)
		case "tcp":
			wg.Add(1)
			go createTCP(v.Port, v.Return, &wg)
		default:
			return fmt.Errorf("unknown type of service to create: %v", v.Type)
		}
	}

	wg.Wait()

	return nil
}

func createUDP(port int, ret interface{}, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	buf := make([]byte, 1024)
	portStr := fmt.Sprintf(":%d", port)
	pc, err := net.ListenPacket("udp", portStr)
	log.Println(err)

	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println("could not read udp packet: ", err)
		}

		log.Println(n, addr, string(buf[:n]))

		if ret != nil {
			val := processUDPReturn(ret)

			_, err := pc.WriteTo(val, addr)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createTCP(port int, ret interface{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Done()
	}

	portStr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Println("error:", err)
		return
	}
	conn.Close()
}

func processUDPReturn(value interface{}) []byte {
	var ret []byte
	switch v := value.(type) {
	case string:
		ret = []byte(v)
	case []interface{}:
		ret = make([]byte, len(v))
		for i := range v {
			if val, ok := v[i].(int); ok {
				ret[i] = byte(val)
			} else {
				log.Fatal("not a byte array")
			}
		}
	case []byte:
		ret = v
	case []int:
		log.Fatal("array in return should be byte magnitude")
	default:
		log.Fatal("that type is not supported:", reflect.TypeOf(v).String())
	}
	return ret
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	return data, err
}

func generateYamlConfig(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(example)

	return nil
}
