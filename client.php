<?php
$sock = socket_create(AF_INET, SOCK_DGRAM, SOL_UDP);
$msg = 'test';

for ($i = 0; $i <= 3; $i++) {
	$message = $msg . $i;
	socket_sendto($sock, $message, strlen($message), 0, '127.0.0.1', 10001);
	//sleep(1);
}

socket_close($sock);