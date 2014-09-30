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


	fmt.Printf("<< --------------- \r\n")
	fmt.Printf("<< Start UDP listener..... \r\n")

	var udp_addr net.UDPAddr
	udp_addr.IP = net.IPv4(0,0,0,0)
	udp_addr.Port = 139
	udp_conn,err := net.ListenUDP("udp",&udp_addr)
	if err != nil {
		fmt.Printf("[Error] %v \r\n",err)
	}

	fmt.Printf("<< UDP listening ....\r\n")

	fmt.Printf("<< --------------- \r\n")
	fmt.Printf("<< Start TCP listener..... \r\n")
	var tcp_addr net.TCPAddr
	tcp_addr.IP = net.IPv4(0,0,0,0)
	tcp_addr.Port = 1153
	tcp_listener,err := net.ListenTCP("tcp",&tcp_addr)
	if err != nil {
		fmt.Printf("[Error] %v \r\n",err)
	}

	fmt.Printf("<< TCP listening ....\r\n")

	go func(listener *net.TCPListener) {
		fmt.Printf("<< Tring to accept ...\r\n")
		for {
			listener.Accept()
			fmt.Printf("<< New connected ...\r\n")
		}
	}(tcp_listener)
	time.Sleep(1 * time.Second)

	fmt.Printf("<< Send dns data to transfer host...\r\n")
	

	for {
		fmt.Printf("\r\n<< Waiting for dns query...\r\n")
		dns_data := make([]byte,2048) 
		length,addr,err := udp_conn.ReadFromUDP(dns_data)
		if err != nil {
			fmt.Printf("<<[Error] %v \r\n",err)
		}

		dns_data  = dns_data[0:length]
		dns_head := dns_data[0:12]
		dns_body := dns_data[12:length-5]

		fmt.Printf("<< [Len]:%v \r\n<< [addr]:%v \r\n",length,addr)
		fmt.Printf("<< [Head]:%v \r\n<< [Body]:%v \r\n",dns_head,dns_body)
		fmt.Printf("<< [Domain]:%v \r\n",parse_domain(dns_body))
		fmt.Printf("\r\n<< Throw in goroutine to handle ....\r\n")

		go func(data []byte,conn *net.UDPConn,udp_addr *net.UDPAddr){
			bytes,err := transfer_dns(data) ; if err != nil {
				fmt.Printf("<<[Error] %v \r\n",err)
			}
			conn.WriteToUDP(bytes,udp_addr)
		}(dns_data,udp_conn,addr)

	} // end for
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





