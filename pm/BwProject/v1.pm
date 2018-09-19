package v1;
use v5.18;
use strict;
use warnings;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
$VERSION = 1.00;
@ISA = qw(Exporter);
@EXPORT_OK = qw/getBasePort/;
%EXPORT_TAGS = ();
@EXPORT = ();

# =============================================================================

my (@vars);
use vars qw/$selfFileSpec/;
BEGIN {
  my @caller = caller(0);
  $selfFileSpec = $caller[1];
  open(IN, $selfFileSpec);
  while (<IN>) {
    push @vars, $1 if /^\s*(\$\w[\w\d]*Def)\s*=/;
  }
  close(IN);
  push @vars, '$defs';

}
use vars @vars;

# =============================================================================

use Data::Dumper;
use Bw;

# =============================================================================

my %basePorts = (
  ssh => 2200,
  http => 8000,
  https => 4400,
  mysql => 3300,
  redis => 6300,
  rabbitmq => 5600,
);

sub getBasePort($) {
  my $key = shift;
  return $basePorts{$key} || die;
}

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
      $cnf->{dockerImageIdFileSpec} = "$ENV{projDir}/docker/dev-$cnf->{projShortcut}.id";
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
    force => {
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
Перед сборкой образа необходимо внести изменения в <ansiFileSpec>${\shortenFileSpec "$ENV{projDir}/$Dockerfile"}<ansi> или выполнить команду с опцией <ansiCmd>--force<ansi> ( <ansiCmd>-f<ansi> )
MSG
  } else {
    use File::Copy::Recursive qw(dircopy);
    dircopy("$ENV{_bwDir}/pm/BwProject/$ENV{bwProjectVersion}/docker", "$ENV{projDir}/docker/.helper") or die;
    my $errorCode = docker({ v => 'all' }, qw/build --pull -t/, $cnf->{dockerImageName}, "$ENV{projDir}/docker");
    if (!$errorCode) {
      $errorCode = docker( { v => 'err' }, qw/inspect --format/, '{{json .Id}}', "$cnf->{dockerImageName}:latest", '>', "$cnf->{dockerImageIdFileSpec}"); exit $errorCode if $errorCode;
      if ( gitIsChangedFile("docker/$cnf->{dockerImageIdFileSpec}", $ENV{projDir}) ) {
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
  options => Hash::Ordered->new(
    noCheck => {
      type => 'bool',
      shortcut => 'n',
      description => sub {
        my $cnf = shift;
        "Не проверять актуальность docker-образа <ansiPrimaryLiteral>$cnf->{dockerImageName}"
      }
    },
    portIncrement => {
      type => 'int',
      min => 0,
      default => 0,
      shortcut => 'i',
      description => sub {
        my $cnf = shift;
        "
        Смещение относительно базовых значений для всех портов
        Полезно для старта второго экземпляра docker-приложения <ansiPrimaryLiteral>${\shortenFileSpec($cnf->{projDir})}<ansi>
        Примечание: второй экземпляр следует запускать из второй копии проекта, которую можно установить командой:
          <ansiCmd>bw p $cnf->{projShortcut} -p <ansiOutline>Абсолютный-путь-к-папке-второй-копии-проекта
        "
      },
    },
    map({
      my $key = $_;
      (
        $key => {
          condition => sub {
            my $cnf = shift;
            $cnf->{mixin}->{scalar getFuncName(3)}->{options}->exists($key)
          },
          type => 'int',
          min => 1024,
          max => 65535,
          description => "<ansiSecondaryLiteral>$key<ansi>-порт по которому будет доступно docker-приложение" . (
            $key ne 'upstream' ? '' : ' для сервиса <ansiPrimaryLiteral>nginx<ansi>'
          ),
        }
      )
    } qw/ssh http https mysql redis rabbitmq upstream/),
    noTestAccessMessage => {
      type => 'bool',
      shortcut => 'm',
      description => 'Не выводить сообщение о проверке доступности docker-приложения',
    },
    forceRecreate => {
      type => 'bool',
      shortcut => 'f',
      description => 'Поднимает docker-приложение с опцией <ansiCmd>--force-recreate<ansi>, передаваемой <ansiCmd>docker-compose up'
    },
    # restart => {
    #   type => 'list',
    #   itemType => 'enum',
    #   enum => sub { my $cnf = shift; [ qw/main nginx/ ] },
    #   shortcut => 'r',
    #   description => 'Останавливает и поднимает указанные сервисы',
    # },
  ),
};
sub cmd_docker_up {
  my ($p, $cnf) = &processParams; $p || return;
  # my $p = &processParams;
  # print Dumper($p);

  # TODO:
  # if [[ -n $https ]]; then
  #   bw_install --silentIfAlreadyInstalled root-cert || { errorCode=$?; break; }
  # fi

  if ( $ENV{OSTYPE} =~ /^linux/ ) {
    my $line='fs.inotify.max_user_watches=524288';
    my $fileSpec='/etc/sysctl.conf';
    my $found = 0;
    # https://perlmaven.com/how-to-grep-a-file-using-perl
    open IN, '<:encoding(UTF-8)', $fileSpec or die;
    while (<IN>) {
      next unless $_ eq $line;
      $found = 1;
      last;
    }
    close IN;
    if (!$found) {
      my $errorCode = execCmd({ v => 'all' }, 'echo', $line, ' | sudo tee -a ', $fileSpec, ' && sudo sysctl -p'); exit $errorCode if $errorCode;
    }
  }

  my %portNameByValue = ();
  foreach my $portName (keys %basePorts) {
    if ( exists $p->{$portName} ) {
      my $portValue = $p->{$portName}->{value};
      if ( exists $portNameByValue{$portValue} ) {
        die ansi 'Err', <<"MSG";
<ansiCmd>$cnf->{cmd}<ansi> обнаружил, что <ansiSecondaryLiteral>$portName<ansi>-порт и <ansiSecondaryLiteral>$portNameByValue{$portValue}<ansi>-порт имеют одинаковое значение <ansiPrimaryLiteral>$portValue
MSG
      }
      $portNameByValue{$portValue} = $portName;
    }
  }

  # TODO:
  #  bw_install --silentIfAlreadyInstalled docker || { errorCode=$?; break; }

  if ($cnf->{dockerImageName} && !$p->{noCheck}->{value}) {
    my ($dockerImageLs, $errorCode) = docker({return => 'stdout', v => 'all'}, qw/image ls/, "$cnf->{dockerImageName}:latest", qw/-q/); exit $errorCode if $errorCode;
    # print Dumper({stdout => $stdout});
    if (!$dockerImageLs) {
      my $errorCode = docker({v => 'all'}, qw/image pull/, "$cnf->{dockerImageName}:latest"); exit $errorCode if $errorCode;
    }
    # my $tstImageIdFileSpec = "/tmp/$cnf->{projShortcut}.image.id";
    my $imageId;
    ($imageId, $errorCode) = docker({ v => 'err', return => 'stdout' }, qw/inspect --format/, '{{json .Id}}', "$cnf->{dockerImageName}:latest"); exit $errorCode if $errorCode;
    my $/ = undef;
    open IN, $cnf->{dockerImageIdFileSpec} or die;
    $etaImageId = <IN>;
    close IN;
    die Dumper({etaImageId => $etaImageId});

    # use File::Compare;
    # if (compare($tstImageIdFileSpec, "$cnf->{dockerImageIdFileSpec}")) {
    #   # gitIsChangedFile
    #   my $needWarn = 0;
    #   if ( gitIsChangedFile('docker/Dockerfile', $cnf->{projDir}) ) {
    #     $needWarn = 1;
    #   } else {
    #     my $errorCode = docker({v => 'all'}, qw/image pull/, "$cnf->{dockerImageName}:latest"); exit $errorCode if $errorCode;
    #     my $errorCode = docker(qw/inspect --format/, '{{json .Id}}', "$cnf->{dockerImageName}:latest", '>', $tstImageIdFileSpec); exit $errorCode if $errorCode;

    #   }
    # }
    # my $errorCode = execCmd('cmp', );
  }
}

# Использование: dip docker up [Опции]
# Описание: Поднимает (docker-compose up) следующие контейнеры:
# Опции
#     --no-check или -n
#         Не проверять актуальность docker-образа bazawinner/dev-dip
#     --ssh значение
#         ssh-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 2208
#     --http значение
#         http-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 8008
#     --https значение
#         https-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 4408
#     --mysql значение
#         mysql-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 3308
#     --redis значение
#         redis-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 6308
#     --rabbitmq значение
#         rabbitmq-порт по которому будет доступно docker-приложение
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 5608
#     --upstream значение
#         upstream-порт по которому будет доступно docker-приложение для сервиса nginx
#         Значение - целое число из диапазона 1024..65535
#         Значение по умолчанию: 3000
#     --no-test-access-message или -m
#         Не выводить сообщение о проверке доступности docker-приложения
#     --proj-dir значение или -p значение
#         Папка проекта
#     --force-recreate или -f
#         Поднимает  с опцией --force-recreate, передаваемой docker-compose up
#     --restart значение или -r значение
#         Останавливает и поднимает указанные сервисы
#         Варианты значения:
#         Опция предназначена для того, чтобы сформировать
#         возможно пустой список значений
#         путем eё многократного использования
#     --help или -? или -h
#         Выводит справку

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
