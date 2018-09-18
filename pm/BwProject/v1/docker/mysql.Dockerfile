
FROM bw

# ==============================================================================
# ========================== MySQL Percona 5.6 =================================

# https://www.percona.com/doc/percona-server/5.6/installation/apt_repo.html
RUN sudo curl -O https://repo.percona.com/apt/percona-release_0.1-6.$(lsb_release -sc)_all.deb
RUN sudo dpkg -i percona-release_0.1-6.$(lsb_release -sc)_all.deb
RUN mkdir /home/dev/mysql
COPY mysql/install/script.exp /home/dev/mysql/
RUN \
	sudo apt-get update && \
  sudo apt-get install -y expect && \
	ls /home/dev/mysql/script.exp && \
  sudo chmod 777 /home/dev/mysql/script.exp && \
	sudo /home/dev/mysql/script.exp && \
true
RUN sudo cp -r /var/lib/mysql /var/lib/_mysql