<?php
$user_id = 1;
$user = json_decode(file_get_contents("http://localhost:8001/api/user/profile?user_id={$user_id}"), true);
?>
<!DOCTYPE html>
<html>
<head><title>Profile</title></head>
<body>
<h2>Welcome, <?php echo htmlspecialchars($user['username']); ?></h2>
<p>Score: <?php echo $user['score']; ?></p>
<p>Test Attempted: <?php echo $user['test_attempted'] ? 'Yes' : 'No'; ?></p>
<a href="test_page.php">Take the test again</a>
</body>
</html>
