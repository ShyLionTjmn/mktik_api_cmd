Usage: ./mktik_api_cmd [OPTIONS] command [command options ...]
	-i	Device IP (mandatory)
	-u	User (default admin, also MKTIK_API_USER environment variable used)
	-p	Password (default "", also MKTIK_API_PASS environment variable used)
	-P	Port (default 8728)
	-d	Debug

Example:
	./mktik_api_cmd -i 10.100.26.160 '/interface/print' '=.proplist=name,type,disabled' '?type=ether'
