package mysql

func (m *MySQL) CodisWrite(cluster_name string, proxy_host string, server_host string) bool {

	m.Create(&CodisInfo{
		ClusterName: cluster_name,
		ProxyHost:   proxy_host,
		ServerHost:  server_host,
	})
	return true
}

func (m *MySQL) CodisUpdate(cluster_name string, proxy_host string, server_host string) bool {
	var codisinfo CodisInfo
	// m.Model(&codisinfo).Where("cluster_name = ?", cluster_name).Update("proxy_host", proxy_host)
	m.Model(&codisinfo).Where("cluster_name = ?", cluster_name).Updates(CodisInfo{ProxyHost: proxy_host, ServerHost: server_host})
	return true
}

func (m *MySQL) CodisExist(cluster_name string) bool {
	var codisinfo CodisInfo
	if m.Where("cluster_name = ?", cluster_name).First(&codisinfo).RecordNotFound() {

		return false
	}
	return true
}

func (m *MySQL) GraphWrite(Client_ip string, Codis_info string) bool {

	m.Create(&CodisGraph{
		ClientIp:  Client_ip,
		CodisInfo: Codis_info,
	})
	return true
}

func (m *MySQL) GraphUpdate(cluster_name string, proxy_host string, server_host string) bool {
	var codisinfo CodisInfo
	// m.Model(&codisinfo).Where("cluster_name = ?", cluster_name).Update("proxy_host", proxy_host)
	m.Model(&codisinfo).Where("cluster_name = ?", cluster_name).Updates(CodisInfo{ProxyHost: proxy_host, ServerHost: server_host})
	return true
}

func (m *MySQL) GraphExist(cluster_name string) bool {
	var codisinfo CodisInfo
	if m.Where("cluster_name = ?", cluster_name).First(&codisinfo).RecordNotFound() {

		return false
	}
	return true
}
