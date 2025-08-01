<?php
$username = $_POST['username'];

$payload = json_encode(['username' => $username]);

$ch = curl_init('http://localhost:8003/api/auth/login');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $payload);
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);

$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($httpCode === 200 && $response) {
    $data = json_decode($response, true);

    if (isset($data['token']) && isset($data['user_id'])) {
        session_start();
        $_SESSION['user_id'] = $data['user_id'];
        $_SESSION['token'] = $data['token'];
        header("Location: test_page.php");
        exit;
    } else {
        echo "Invalid response from auth service.";
    }
} else {
    echo "Login failed.";
}
?>
