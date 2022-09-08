#include <stdio.h>
#include <sys/socket.h>
#include <string.h>
#include <netinet/in.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <pthread.h>

#define PORT     13337

typedef struct ReceiveLoopArguments {
    int sock;
    char *local_ip;
} ReceiveLoopArguments;

void receive_callback(char buf[65507], struct sockaddr_in addr, char *local_ip) {
    char ip[INET_ADDRSTRLEN];
    uint16_t port;

    inet_ntop(AF_INET, &(addr.sin_addr), ip, INET_ADDRSTRLEN);
    port = htons (addr.sin_port);

    if (strcmp(ip, local_ip) == 0 && port == PORT) {
        return;
    }

    printf("%s:%d\t%s\n", ip, port, buf);
}

_Noreturn void receive_loop(ReceiveLoopArguments *arguments) {
    char buf[65507];

    struct sockaddr_in addr;
    socklen_t addr_size = sizeof(addr);

    for (;;) {
        memset(&buf, 0, sizeof(buf));

        if (recvfrom(arguments->sock, buf, sizeof(buf), 0, (struct sockaddr *) &addr, &addr_size) < 0) {
            perror("failed to receive from socket");
            exit(1);
        }

        receive_callback(buf, addr, arguments->local_ip);
    }
}

int main() {
    char *hostname = getenv("HOSTNAME");
    char *local_ip = getenv("LOCAL_IP");
    char *broadcast_ip = getenv("BROADCAST_IP");
    int sock;
    struct sockaddr_in local_addr, broadcast_addr;
    char buf[65507];
    ReceiveLoopArguments arguments;
    pthread_t receive_loop_thread;

    if (!hostname) {
        perror("SESSION_UUID env var missing");
        return 1;
    }

    if (!local_ip) {
        perror("LOCAL_IP env var missing");
        return 1;
    }

    if (!broadcast_ip) {
        perror("BROADCAST_IP env var missing");
        return 1;
    }

    if ((sock = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("failed to create socket");
        return 1;
    }

    int enable_broadcast = 1;
    if (setsockopt(sock, SOL_SOCKET, SO_BROADCAST, &enable_broadcast, sizeof(enable_broadcast)) < 0) {
        perror("failed to set broadcast on socket");
        return 1;
    }

    memset(&local_addr, 0, sizeof(local_addr));
    local_addr.sin_family = AF_INET;
    local_addr.sin_addr.s_addr = INADDR_ANY;
    local_addr.sin_port = htons(PORT);

    memset(&broadcast_addr, 0, sizeof(broadcast_addr));
    broadcast_addr.sin_family = AF_INET;
    broadcast_addr.sin_addr.s_addr = inet_addr(broadcast_ip);
    broadcast_addr.sin_port = htons(PORT);

    if (bind(sock, (const struct sockaddr *) &local_addr, sizeof(local_addr)) < 0) {
        perror("failed to bind socket");
        return 1;
    }

    memset(&arguments, 0, sizeof(arguments));
    arguments.sock = sock;
    arguments.local_ip = local_ip;

    pthread_create(&receive_loop_thread, NULL, (void *) receive_loop, (void *) &arguments);

    sprintf(buf, "Hello world from C @ %s", hostname);

    for (;;) {
        if (sendto(sock, buf, strlen(buf), 0, (const struct sockaddr *) &broadcast_addr, sizeof(broadcast_addr)) < 0) {
            perror("failed to send to socket");
            return 1;
        }

        sleep(1);
    }
}
