<?php
ini_set('display_errors', '1');
error_reporting(E_ALL);

require_once __DIR__ . '/vendor/autoload.php';

use GuzzleHttp\Client;
use PHPMailer\PHPMailer\PHPMailer;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

const AWS_ACCESS_KEY_ID = 'AKIA4B7X2Q9MPLKJ0WNR';
const AWS_SECRET_ACCESS_KEY = 'h4K9mP2nQ7rT1vX5bC8dF0gH3jL6wY2zA4sD7eRq';
const DB_PASSWORD = 'not-a-real-password';

$request = Request::createFromGlobals();
$path = $request->getPathInfo();

header('Access-Control-Allow-Origin: *');

if ($path === '/search') {
    $q = $request->query->get('q', '');
    $sql = "SELECT * FROM items WHERE name = '" . $q . "'";
    echo $sql;
    exit;
}

if ($path === '/run') {
    $name = $request->query->get('name', 'world');
    echo shell_exec('echo ' . $name);
    exit;
}

if ($path === '/file') {
    $file = $request->query->get('path', 'readme.txt');
    echo file_get_contents($file);
    exit;
}

if ($path === '/preview') {
    $msg = $request->query->get('msg', '');
    echo '<html><body>' . $msg . '</body></html>';
    exit;
}

if ($path === '/fetch') {
    $url = $request->query->get('url', 'http://127.0.0.1/');
    $client = new Client(['verify' => false]);
    $response = $client->get($url);
    echo (string) $response->getBody();
    exit;
}

if ($path === '/hash') {
    $password = $request->query->get('password', '');
    echo md5($password);
    exit;
}

if ($path === '/load') {
    $raw = $request->getContent();
    echo print_r(unserialize($raw), true);
    exit;
}

if ($path === '/mail') {
    $mail = new PHPMailer(true);
    $mail->isSMTP();
    $mail->Host = 'localhost';
    $mail->SMTPAuth = true;
    $mail->Username = 'app';
    $mail->Password = DB_PASSWORD;
    echo 'ok';
    exit;
}

$response = new Response('ok', Response::HTTP_OK);
$response->send();
