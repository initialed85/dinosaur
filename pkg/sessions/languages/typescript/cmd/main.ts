import * as dgram from 'dgram';

const PORT = 13337;

function receiveCallback(data: string, ip: string, port: number, localIp: string) {
    if (ip === localIp && port === PORT) {
        return;
    }

    console.log(`${ip}:${port}\t${data}`);
}

function main() {
    const hostname = process.env.HOSTNAME;
    const localIp = process.env.LOCAL_IP;
    const broadcastIp = process.env.BROADCAST_IP;

    const sock = dgram.createSocket('udp4');

    sock.on('error', err => {
        console.log(`error: ${err.stack}`);
        sock.close();
        return;
    });

    sock.on('message', (msg, rinfo) => {
        receiveCallback(msg.toString(), rinfo.address, rinfo.port, localIp);
    });

    sock.on('listening', () => {
        sock.setBroadcast(true);

        setInterval(() => {
            sock.send(`Hello from TypeScript @ ${hostname}`, PORT, broadcastIp);
        }, 1_000);
    });

    sock.bind(PORT, '0.0.0.0');
}

main();
