<?php
$sock = socket_create(AF_INET, SOCK_DGRAM, SOL_UDP);

for ($i = 0; $i <= 1000; $i++) {
	$value = [
		'ip' => '127.0.0.1'
	];
	$message = json_encode([
		'name' => 'record_' . $i,
		'group' => 'test',
		'data' => bin2hex(gzcompress(serialize($value)))
	]);
	socket_sendto($sock, $message, strlen($message), 0, '127.0.0.1', 10001);
	echo 'send message ' . $message . PHP_EOL;
	sleep(1);
}

socket_close($sock);