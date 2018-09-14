package v1;
use v5.18;
use strict;
use warnings;

my (@vars);
use vars qw/$selfFileSpec/;
BEGIN {
  use BwCore;

  # use Exporter;
  # use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
  # $VERSION = 1.00;
  # @ISA = qw(Exporter);
  # @EXPORT = ();
  # @EXPORT_OK = ();
  # %EXPORT_TAGS = ();

  my @caller = caller(0);
  $selfFileSpec = $caller[1];
  open(IN, $selfFileSpec);
  while (<IN>) {
    push @vars, $1 if /^\s*(\$\w[\w\d]*Def)\s*=/;
  }
  close(IN);
  push @vars, '$defs';

  # @EXPORT = qw/cmd/;
}
use vars @vars;

# =============================================================================

sub preprocessCnf {
  my $entry = shift;
  my $cnf = shift;

  if ( $entry eq 'cmd' ) {
    if (!$cnf->{noDockerBuild}) {
      my $dockerImageName = $cnf->{dockerImageName};
      my $domain = 'bazawinner';
      if (!$dockerImageName) {
        $cnf->{dockerImageName} = "$domain/dev-$cnf->{projShortcut}"
      } elsif ( index($dockerImageName, '/') == -1 ) {
        $cnf->{dockerImageName} = "$domain/$cnf->{dockerImageName}"
      }
    }
  }

  $cnf;
}

# =============================================================================

$cmdDef = {
  description => sub {
    my $cnf = shift;
    "Базовая утилита проекта <ansiPrimaryLiteral>${\shortenFileSpec($cnf->{projDir})}"
  },
  isWrapper => 1,
};
sub cmd {
  &processParams;
}

# =============================================================================

# Команда - один из следующих вариантов: api docker mysql mysqldump prepare-local self-test st test update worker
#     dip api
#         (Пере)Запускает api
#     dip docker
#         docker-операции
#     dip mysql
#         Запускает mysql-клиент
#     dip mysqldump
#         Запускает mysqldump
#     dip prepare-local
#         Производит подготовку локальных инструментов (баз данных и т.п)
#     dip self-test
#     dip st
#         Самопроверка
#     dip test
#         Прогоняет тест(ы)
#     dip update
#         Обновляет команду dip
#     dip worker
#         (Пере)Запускает worker
# Подробнее см.
#     dip Команда --help
#     dip Команда -?
#     dip Команда -h
# HELP

# =============================================================================

$cmd_dockerDef = {
  'description' => 'docker-операции',
  isWrapper => 1,
};
sub cmd_docker {
  &processParams;
}

# =============================================================================

$cmd_docker_buildDef = {
  condition => sub {
    my $cnf = shift;
    !$cnf->{noDockerBuild};
  },
  description => sub {
    my $cnf = shift;
    "Собирает docker-образ <ansiPrimaryLiteral>$cnf->{dockerImageName}";
  },
  options => Hash::Ordered->new(
    'force' => {
      type => 'bool',
      shortcut => 'f',
      description => 'Невзирая на возможное отсутствие изменений в <ansiFileSpec>docker/Dockerfile'
    },
  ),
};
sub cmd_docker_build {
  my ($p, $cnf) = &processParams; $p || return;
  my $Dockerfile = 'docker/Dockerfile';
  if ( !$p->{force}->{value} && !gitIsChangedFile($Dockerfile, $ENV{projDir}) ) {
    print ansi 'Warn', <<"MSG";
Перед сборкой образа необходимо внести изменения в <ansiFileSpec>${\shortenFileSpec "$ENV{projDir}/$Dockerfile"}<ansi> или выполнить команду с опцией <ansiCmd>--force<ansi> ( <ansiCmd>-f<ansi> )"
MSG
  } else {
    if (docker(qw/build --pull -t/, $cnf->{dockerImageName}, "$ENV{projDir}/docker")) {
      my $dockerImageIdFileName="dev-$cnf->{projShortcut}.id";
      docker(qw/inspect --format "{{json .Id}}"/, "$cnf->{dockerImageName}:latest", '>', "$dockerImageIdFileName");
      if ( gitIsChangedFile("docker/$dockerImageIdFileName", $ENV{projDir}) ) {
        print ansi 'Warn', <<"MSG";
Обновлен docker-образ <ansiPrimaryLiteral>$cnf->{dockerImageName}<ansi>
Не забудьте поместить его в docker-репозиторий командой
    <ansiCmd>$cnf->{projShortcut} docker push
MSG
      }
    }
  }
}

# Использование: dip docker build [Опции]
# Описание: Собирает docker-образ bazawinner/dev-dip
# Опции
#     --force
#         Невзирая на возможное отсутствие изменений в docker/Dockerfile
#     --proj-dir значение или -p значение
#         Папка проекта
#     --help или -? или -h
#         Выводит справку

# =============================================================================

$cmd_docker_pushDef = {
  condition => sub {
    my $cnf = shift;
    !$cnf->{noDockerBuild};
  },
  description => sub {
    my $cnf = shift;
    "Пушит на <ansiUrl>https://hub.docker.com/<ansi> docker-образ <ansiPrimaryLiteral>$cnf->{dockerImageName}";
  },
};
sub cmd_docker_push {
  my $p = &processParams;
  print Dumper($p);
}

# =============================================================================

$cmd_docker_upDef = {
  description => sub {
    my $cnf = shift;
    "Поднимает командой <ansiCmd>docker-compose up<ansi> docker-приложение <ansiPrimaryLiteral>${\shortenFileSpec($cnf->{projDir})}";
  },
};
sub cmd_docker_up {
  my $p = &processParams;
  print Dumper($p);
}

# =============================================================================

$cmd_docker_shellDef = {
  description => sub {
    my $cnf = shift;
    "Запускает shell в docker-контейнере";
  },
};
sub cmd_docker_shell {
  my $p = &processParams;
  print Dumper($p);
}
# Использование: dip docker [Опции] Команда
# Описание: docker-операции
# Команда - один из следующих вариантов: build down push shell up
#     dip docker build
#         Собирает docker-образ bazawinner/dev-dip
#     dip docker down
#         Останавливает (docker-compose down) следующие контейнеры:
#     dip docker push
#         Push-ит docker-образ bazawinner/dev-dip
#     dip docker shell
#         Запускает bash в docker-контейнере
#     dip docker up
#         Поднимает (docker-compose up) следующие контейнеры:
# Подробнее см.
#     dip docker Команда --help
#     dip docker Команда -?
#     dip docker Команда -h
# Опции
#     --help или -? или -h
#         Выводит справку

# =============================================================================

$cmd_selfTestDef = {
  'description' => 'Самопроверка',
  'shortcut' => 'st',
};
sub cmd_selfTest {
  my $p = &processParams;
  print Dumper($p);
  # my $cnf = shift;
  # print Dumper(getFuncName(), $cnf, @_);
  # print "selfTest\n"
}

$defs = preprocessDefs;
