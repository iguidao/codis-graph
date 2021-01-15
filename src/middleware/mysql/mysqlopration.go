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

func (m *MySQL) CodisProxyLike(Proxy_ip string) (string, uint) {
	var codisinfo CodisInfo
	m.Where("proxy_host LIKE ?", "%,"+Proxy_ip+",%").First(&codisinfo)

	return codisinfo.ClusterName, codisinfo.ID
}
func (m *MySQL) GraphWrite(client_ip, cluster_name string, cluster_id uint) bool {
	m.Create(&CodisGraph{
		ClientIp:    client_ip,
		ClName:      cluster_name,
		CodisInfoID: cluster_id,
		// CodisInfo: CodisInfo{ClusterName: cluster_name},
	})
	return true
}

func (m *MySQL) GraphUpdate(client_ip, cluster_name string, cluster_id uint) bool {
	var codisgraph CodisGraph
	m.Model(&codisgraph).Where("client_ip = ?", client_ip).Updates(CodisGraph{ClName: cluster_name, CodisInfoID: cluster_id})
	return true
}
func (m *MySQL) GraphDel(client_ip string) bool {
	var codisgraph CodisGraph
	m.Where("client_ip = ?", client_ip).Delete(&codisgraph)
	return true
}

func (m *MySQL) GraphExist(client_ip string) bool {
	var codisgraph CodisGraph
	if m.Where("client_ip = ?", client_ip).First(&codisgraph).RecordNotFound() {
		return false
	}
	return true
}

func (m *MySQL) GraphGetAll() []CodisGraph {
	var codisgraph []CodisGraph
	m.Preload("CodisInfo").Find(&codisgraph)
	return codisgraph
}
