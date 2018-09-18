package Bw;
use v5.18;
use strict;
use warnings;
use Exporter;
use vars qw($VERSION @ISA @EXPORT @EXPORT_OK %EXPORT_TAGS);
$VERSION = 1.00;
@ISA = qw(Exporter);
@EXPORT_OK = ();
%EXPORT_TAGS = ();
@EXPORT = ();

# BEGIN {
  use File::Find qw/find/;
  my $carpAlwaysIsInstalled;
  no warnings 'File::Find';
  find { wanted => sub { $carpAlwaysIsInstalled = 1 if /\/Carp\/Always(?:\.pm)?$/ }, no_chdir => 1 }, @INC;
  if ( $carpAlwaysIsInstalled ) {
    require Carp::Always;# https://metacpan.org/pod/Carp::Always
    Carp::Always->import;
  }
# }

use Hash::Ordered; # https://metacpan.org/pod/Hash::Ordered

use Data::Dumper;
push @EXPORT, qw/Dumper/;

use BwCore;
push @EXPORT, @BwCore::EXPORT;

use BwAnsi;
push @EXPORT, @BwAnsi::EXPORT;

use BwParams;
push @EXPORT, @BwParams::EXPORT;

1;