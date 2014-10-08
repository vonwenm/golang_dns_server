package main


import (
	"fmt"
	"net"
	"time"
)


func main() {
	//fmt.Printf("<< Server started ! \r\n")
	//addrs,err := net.LookupHost("www.qq.com")
	//if err != nil {
	//	fmt.Printf("[Error] %v \r\n",err)
	//	return
	//}
	//fmt.Printf("<< Addrs:%v \r\n",addrs)

	fmt.Printf("<< Start TCP listener..... \r\n")
	var tcp_addr net.TCPAddr
	tcp_addr.IP = net.IPv4(0,0,0,0)
	tcp_addr.Port = 5553
	tcp_listener,err := net.ListenTCP("tcp",&tcp_addr)
	if err != nil {
		fmt.Printf("[Error] %v \r\n",err)
	}

	fmt.Printf("<< TCP listening ....\r\n")

	go func(listener *net.TCPListener) {
		fmt.Printf("<< Tring to accept ...\r\n")
		for {
			socket_conn,err := listener.Accept()
			fmt.Printf("<< New connected conn:%v error:%v...\r\n",socket_conn,err)

			go func(socket net.Conn) {
				data := make([]byte,1024)
				n,err := socket.Read(data) ; if err != nil {
					fmt.Printf("[Error] %v \r\n",err)
				}
				fmt.Printf("<< [Len]:%v \r\n<< [Data]:%v \r\n",n,data[0:n])
				fmt.Printf("<< Waiting for transfer server repling....\r\n")
				transfer_data,err := transfer_dns(data[0:n]) ; if err != nil {
					fmt.Printf("[Error] %v \r\n",err)
					return 
				}

				n,err = socket.Write(transfer_data) ; if err != nil {
					fmt.Printf("[Error] %v \r\n",err)
					return
				}
				fmt.Printf("<< Translate complete....\r\n")
			}(socket_conn)
			
		}
	}(tcp_listener)

	for {
		time.Sleep(1 * time.Second)	
	}
	
}


func transfer_dns(query_data []byte) ([]byte,error) {
		var (
			conn net.Conn
			err  error
			length int 
		)
		b := make([]byte,2048)

		if conn, err = net.Dial("udp", "8.8.8.8:53"); err != nil {
			fmt.Println(err.Error())
			return nil,err
		}

		if length,err = conn.Write(query_data); err != nil {
			fmt.Println(err.Error())
			return nil,err
		}

		defer conn.Close()

		conn.SetDeadline( time.Now().Add(5 * time.Second) )

		fmt.Printf("<< Send %d bytes to transfered server ... \r\n",length)
		
		length,err = conn.Read(b) ; if err != nil {
			fmt.Println(err.Error())
			return nil,err
		}


		transfered_dns_data := b[0:length]
		transfered_dns_head := b[0:12]
		transfered_dns_body := b[12:length]

		fmt.Printf("<< [Head] : %v\r\n",transfered_dns_head)
		fmt.Printf("<< [Body] : %v\r\n",transfered_dns_body)

		return transfered_dns_data,nil

}


func parse_domain(b []byte) string {
	domain := ""
	start := 0
	end := 0
	section_lenth := 0
	for  {
		section_lenth = int(b[0]) 
		start = start + 1
		end = start  + section_lenth 

		domain = domain + fmt.Sprintf("%s.",b[start:end])

		b = b[end:len(b)]
		start = 0
		
		if len(b) == 0 {
			domain = domain[0:len(domain)-1]
			break
		}
	}
	return domain
}





