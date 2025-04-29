<?php
function api_post($url, $data) {
    $opts = ['http' => [
        'method'  => 'POST',
        'header'  => "Content-Type: application/json",
        'content' => json_encode($data),
    ]];
    $context = stream_context_create($opts);
    return file_get_contents($url, false, $context);
}

$user_id = (int)$_POST['user_id'];
$answers = [];

foreach ($_POST as $key => $value) {
    if (strpos($key, 'question_') === 0) {
        $qid = (int)str_replace('question_', '', $key);
        $answers[$qid] = (int)$value;
    }
}

api_post('http://localhost:8002/api/test/submit', [
    'user_id' => $user_id,
    'answers' => $answers
]);

header("Location: profile.php");
exit;
?>
