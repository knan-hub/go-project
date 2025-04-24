from flask import Flask, request, jsonify
import subprocess

app = Flask(__name__)

@app.route('/convert', methods=['POST'])
def convert():
    try:
        # 从请求中获取文件、from 和 to 参数
        file = request.files['file']
        file_buffer = file.read()
        from_format = request.form.get('from')
        to_format = request.form.get('to')

        # 校验 from 和 to 参数
        if not from_format or not to_format:
            return jsonify({'success': False, 'message': 'Missing from or to format'}), 400

        # 转换文件
        convert = convert_file(file_buffer, from_format, to_format)
        return convert, 200, {'Content-Type': 'text/plain; charset=utf-8'}

    except Exception as err:
        app.logger.error(err)
        return jsonify({'success': False, 'error': str(err)}), 500

def convert_file(file_buffer, from_format, to_format):
    # 使用 subprocess 调用 pandoc
    process = subprocess.Popen(
        ['pandoc', '-f', from_format, '-t', to_format],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    
    stdout, stderr = process.communicate(input=file_buffer)

    if process.returncode != 0:
        raise Exception(stderr.decode('utf-8'))
    
    return stdout.decode('utf-8')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80)
