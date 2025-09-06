<!DOCTYPE html>
<html>

<head>
    <title>&#129699;</title>
    <link rel="stylesheet" href="/static/index.css">
</head>

<body>
    <p>{{.Message}}</p>

    <ul>
        <li><a class="active">admin</a></li>
        <li><a href="https://github.com/timokats/emmer">git</a></li>
        <li><a href="#contact">Contact</a></li>
    </ul>

    <button onclick="request('/api/ping')">Ping Server</button>
    <textarea id="result"> </textarea>

    <script>
        document.getElementById('result').value = 'send request...';

        async function request(path) {
            const response = await fetch(path);
            const text = await response.text();
            document.getElementById('result').value = text;
        }

        const textarea = document.getElementById('result');

        textarea.addEventListener('keydown', function(e) {
            if (e.key === 'Tab') {
                e.preventDefault();
                const start = this.selectionStart;
                const end = this.selectionEnd;
                this.value = this.value.substring(0, start) + "\t" + this.value.substring(end);
                this.selectionStart = this.selectionEnd = start + 1;
            }
        })
    </script>
</body>

</html>