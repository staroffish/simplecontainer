#!/usr/bin/perl

if(scalar(@ARGV) < 1)
{
    print("usage: install.pl INSTALL_DIR\n");
    exit(-1);
}

if($ENV{USER} ne "root")
{
    print("This tool must execute by root\n");
    exit(-1);
}

my $installPath = $ARGV[0];

if($installPath =~ /^\./)
{
    print("Install path must be a absolute path\n");
    exit(-1);
}

unless(mkdir($installPath, 0755))
{
    printf("mkdir $installPath error. path=$!\n");
    exit(-1);
}

if(! -d $installPath || ! -x $installPath || ! -w $installPath) 
{
    printf("Permisson error or not a directory or not exists. path=$installPath\n");
    exit(-1);
}

my $mntPath = $installPath . "/mnt";
my $writelayPath = $installPath . "/writelayer";
my $imagePath = $installPath . "/images";
my $logPath = $installPath . "/log";
my $cInfoPath = $installPath . "/containerInfo";

if(!mkdir($mntPath, 0755))
{
    printf("mkdir $mntPath error. path=$!\n");
    exit(-1);
}
if(!mkdir($writelayPath, 0755))
{
    printf("mkdir $writelayPath error. path=$!\n");
    exit(-1);
}
if(!mkdir($imagePath, 0755))
{
    printf("mkdir $imagePath error. path=$!\n");
    exit(-1);
}
if(!mkdir($logPath, 0755))
{
    printf("mkdir $logPath error. path=$!\n");
    exit(-1);
}
if(!mkdir($cInfoPath, 0755))
{
    printf("mkdir $cInfoPath error. path=$!\n");
    exit(-1);
}

my $config = sprintf("{\n\"mntpath\":\"%s\",\n\"writelaypath\":\"%s\",\n\"imagepath\":\"%s\",\n\"logpath\":\"%s\",\n\"cInfopath\":\"%s\",\n\"loglevel\":\"info\"\n}\n", 
                        $mntPath, $writelayPath, $imagePath, $logPath, $cInfoPath);

unless(open(FD, ">", "/etc/sc.json"))
{
    printf("open /etc/sc.json error. path=$!\n");
    exit(-1);
}

print FD ($config);
close(FD);