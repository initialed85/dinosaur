use std::{env, io, net, str, sync, thread, time};

const PORT: u16 = 13337;
const ONE_SECOND: time::Duration = time::Duration::from_secs(1);

fn receive_callback(src: net::SocketAddr, data: &str, local_ip: String) {
    if src.ip().to_string() == local_ip && src.port() == PORT {
        return;
    }

    println!("{}\t{:?}", src, data);
}

fn receive_loop(socket: sync::Arc<net::UdpSocket>, local_ip: String) {
    let mut buf = [0; 65507];

    loop {
        let result = socket.recv_from(&mut buf);

        if result.is_err() && result.as_ref().unwrap_err().kind() == io::ErrorKind::WouldBlock {
            continue;
        }

        let (n, src) = result.expect("failed to receive from socket");
        let data = str::from_utf8(&buf[..n]).expect("failed to convert bytes to string");

        receive_callback(src, data, local_ip.clone());
    }
}

fn main() -> Result<(), io::Error> {
    let hostname: String = env::var("HOSTNAME").expect("HOSTNAME env var missing");
    let local_ip: String = env::var("LOCAL_IP").expect("LOCAL_IP env var missing");
    let broadcast_ip: String = env::var("BROADCAST_IP").expect("BROADCAST_IP env var missing");

    let local_addr = format!("0.0.0.0:{}", PORT);

    let socket = sync::Arc::new(net::UdpSocket::bind(local_addr).expect("failed to bind socket"));

    socket.set_read_timeout(Option::from(ONE_SECOND)).expect("failed to set read timeout on socket");
    socket.set_broadcast(true).expect("failed to set broadcast on socket");

    let arc_socket = sync::Arc::clone(&socket);

    thread::spawn(move || receive_loop(arc_socket, local_ip));

    loop {
        socket.send_to(
            format!("Hello world from Rust @ {}", hostname).as_bytes(),
            format!("{}:{}", broadcast_ip, PORT),
        ).expect("failed to send to socket");
        thread::sleep(ONE_SECOND);
    }
}
