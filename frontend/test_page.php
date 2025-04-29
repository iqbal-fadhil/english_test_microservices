<?php
$user_id = 1;
$questions = json_decode(file_get_contents("http://localhost:8002/api/test/questions"), true);
?>
<!DOCTYPE html>
<html>
<head><title>TOEFL Test</title></head>
<body>
<form method="POST" action="submit_test.php">
    <?php foreach ($questions as $q): ?>
        <p><strong><?php echo $q['text']; ?></strong></p>
        <?php foreach ($q['choices'] as $c): ?>
            <label>
                <input type="radio" name="question_<?php echo $q['id']; ?>" value="<?php echo $c['id']; ?>" required>
                <?php echo $c['text']; ?>
            </label><br>
        <?php endforeach; ?>
        <hr>
    <?php endforeach; ?>
    <input type="hidden" name="user_id" value="<?php echo $user_id; ?>">
    <button type="submit">Submit Test</button>
</form>
</body>
</html>
