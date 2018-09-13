package ProjectSupport;
use strict;
use warnings;

# =============================================================================

use lib "$ENV{_bwDir}";
use bw;
use Data::Dumper;

use vars qw/$cmdDef/;
$cmdDef = {
  description => sub {
    my $arg = shift;
    "Базовая утилита проекта ${ansiPrimaryLiteral}${\shortenFileSpec($arg->{projDef}->{dir})}${ansiReset}"
  },
  isWrapper => 1,
};
sub cmd {
	my $projDef = shift;
  my $bwProjShortcut = $projDef->{shortcut};
  my $param = shift;
  my $subCommands = getSubCommands($projDef->{pmFileSpec});
	if ($param and ( $param eq '-?' || $param eq '-h' || $param eq '--help')) {
    print <<"HELP";
${ansiHeader}Использование:${ansiReset} ${ansiCmd}$bwProjShortcut${ansiReset} [${ansiOutline}Опция${ansiReset}] ${ansiOutline}Команда${ansiReset}
${ansiHeader}Описание:${ansiReset} ${\getDescription($cmdDef, { projDef => $projDef })}
Команда - один из следующих вариантов: ${ansiSecondaryLiteral}$subCommands->{cmdNameAndShortcuts}${ansiReset}
HELP
    foreach my $cmdName ($subCommands->{byName}->keys) {
      my $cmd = $subCommands->{byName}->get($cmdName);
      my @cmdNameAndShortcuts;
      push @cmdNameAndShortcuts, @{$cmd->{def}->{shortcuts}} if $cmd->{def}->{shortcuts};
      push @cmdNameAndShortcuts, $cmdName;
      foreach my $cmdNameOrShortcut (@cmdNameAndShortcuts) {
        print <<"HELP";
    ${ansiCmd}$bwProjShortcut $cmdNameOrShortcut${ansiReset}
HELP
      }
      print <<"HELP";
        ${\getDescription($cmd->{def}, { projDef => $projDef })}
HELP
    }
    print <<"HELP";
Подробнее см.
    ${ansiCmd}$bwProjShortcut ${ansiOutline}Команда ${ansiCmd}--help${ansiReset}
    ${ansiCmd}$bwProjShortcut ${ansiOutline}Команда ${ansiCmd}-?${ansiReset}
    ${ansiCmd}$bwProjShortcut ${ansiOutline}Команда ${ansiCmd}-h${ansiReset}
Опция
    ${ansiCmd}--help${ansiReset} или ${ansiCmd}-?${ansiReset} или ${ansiCmd}-h${ansiReset}
        Выводит справку
HELP
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
	} else {
    my $subCommand = ! $param ? undef : $subCommands->{byNameOrShortcut}->get($param);
    if ( $subCommand ) {
      &{\&{$subCommand->{funcName}}}(); # https://stackoverflow.com/questions/1915616/how-can-i-elegantly-call-a-perl-subroutine-whose-name-is-held-in-a-variable
    } else {
      my $prefix = !$param ? 'в качесте первого аргумента' : "вместо ${ansiPrimaryLiteral}$param${ansiErr}";
      die "${ansiErr}ERR: ${ansiCmd}$bwProjShortcut${ansiReset}${ansiErr} ${prefix} ожидает одну из следующих команд ${ansiSecondaryLiteral}$subCommands->{cmdNameAndShortcuts}${ansiReset}";
    }
  }
}

use vars qw/$cmd_dockerDef/;
$cmd_dockerDef = {
  'description' => 'docker-операции',
};
sub cmd_docker {
  print Dumper(getFuncName());
}

use vars qw/$cmd_selfTestDef/;
$cmd_selfTestDef = {
  'description' => 'Самопроверка',
  'shortcuts' => [ 'st' ],
};
sub cmd_selfTest {
  print "selfTest\n"
}

1;