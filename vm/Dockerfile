FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Basic tools
RUN echo 'Acquire::Retries "3";' > /etc/apt/apt.conf.d/80-retries && \
    echo 'Acquire::http::Timeout "10";' >> /etc/apt/apt.conf.d/80-retries && \
    sed -i 's|http://archive.ubuntu.com/ubuntu|http://us.archive.ubuntu.com/ubuntu|g' /etc/apt/sources.list && \
    apt-get clean && apt-get update --fix-missing

RUN apt-get update && apt-get install -y --fix-missing \
    bash \
    curl \
    git \
    supervisor \
    sudo \
    wget \
    gpg \
    build-essential \
    binutils \
    gdb \
    file \
    unzip \
    vim \
    nano \
    htop \
    net-tools \
    python3 \
    python3-pip \
    python3-venv \
 && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:mozillateam/ppa && \
    echo 'Package: *\nPin: release o=LP-PPA-mozillateam\nPin-Priority: 1001' \
        > /etc/apt/preferences.d/mozilla-firefox && \
    apt-get update && \
    apt-get install -y firefox && \
    rm -rf /var/lib/apt/lists/*

# Desktop
RUN apt-get update && apt-get install -y --fix-missing \
    xfce4 \
    xfce4-terminal \
    xfce4-goodies \
    dbus-x11 \
    xauth \
    xfonts-base \
    x11-xserver-utils \
    xdg-utils \
    fonts-dejavu \
    xfce4-notifyd \
    notification-daemon \
 && rm -rf /var/lib/apt/lists/*

# VNC
RUN apt-get update && apt-get install -y --fix-missing \
    tigervnc-standalone-server \
    tigervnc-common \
 && rm -rf /var/lib/apt/lists/*

# Utilities
RUN apt-get update && apt-get install -y --fix-missing \
    xterm \
 && rm -rf /var/lib/apt/lists/*

# Create user
RUN useradd -m -s /bin/bash vm && \
    echo "vm:vm" | chpasswd && \
    echo "vm ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# VNC xstartup
COPY xstartup /home/vm/.vnc/xstartup
RUN chown vm:vm /home/vm/.vnc/xstartup && \
    chmod +x /home/vm/.vnc/xstartup && \
    mkdir -p /tmp/.X11-unix && chmod 1777 /tmp/.X11-unix

# Supervisor config
COPY supervisord.conf /etc/supervisor/conf.d/

ENV HOME=/home/vm
ENV DISPLAY=:1
ENV USER=vm
ENV SHELL=/bin/bash
ENV XDG_RUNTIME_DIR=/tmp/runtime-vm

RUN mkdir -p /tmp/runtime-vm && \
    chown vm:vm /tmp/runtime-vm && \
    chmod 755 /tmp/runtime-vm

# Remove xfce4-panel's systray plugin to stop crash popups
RUN mkdir -p /home/vm/.config/xfce4/xfconf/xfce-perchannel-xml && \
    echo '<?xml version="1.0" encoding="UTF-8"?> \
<channel name="xfce4-panel" version="1.0"> \
  <property name="plugins" type="array"> \
    <!-- no systray plugin --> \
  </property> \
</channel>' \
    > /home/vm/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-panel.xml && \
    chown -R vm:vm /home/vm/.config

# Disable unwanted autostart daemons
RUN mkdir -p /home/vm/.config/autostart && \
    for f in \
      nm-applet.desktop \
      xfce-polkit.desktop \
      print-applet.desktop \
      xfce4-power-manager.desktop \
    ; do \
      echo "[Desktop Entry]\nType=Application\nName=$f\nHidden=true" \
        > /home/vm/.config/autostart/$f ; \
    done && \
    chown -R vm:vm /home/vm/.config

EXPOSE 5901

WORKDIR /home/vm

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
