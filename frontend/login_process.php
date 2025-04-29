<?php
$username = $_POST['username'];

$response = file_get_contents("http://localhost:8001/api/user/find?username=" . urlencode($username));
$user = json_decode($response, true);

if ($user && isset($user['id'])) {
    session_start();
    $_SESSION['user_id'] = $user['id'];
    header("Location: test_page.php");
    exit;
} else {
    echo "Invalid login.";
}
?>
