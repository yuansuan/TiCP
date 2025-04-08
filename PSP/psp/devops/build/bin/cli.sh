#!/usr/bin/env bash

set -x

usage_help()
{
    echo "Usage:"
    echo "      gen [j2g/y2g/p2g] [file or dir(include .json or.yml or .proto)]"
    echo ""
    exit 1
}

if [ -z "$*" ]; then
	funcHelp
fi

build_tool_path=devops/build/build_tool

curpath=`pwd`
len_path=${#ROOTPATH}
GEN_COMMD=node

YAML_SCRIPT=build_tool/yaml2go/yaml-to-go.js
JSON_SCRIPT=build_tool/json2go/json-to-go.js
DIR_GEN=generated
MKDIR=mkdir
PARSER_SCRIPT=build_tool/parser/parser.go

mac_os=Darwin

basepath=$(cd `dirname $0`; cd ../ ; pwd)

gen(){
	if [ X"$1" = X"y2g" ]; then
		$GEN_COMMD $basepath/$YAML_SCRIPT $2 $3
	elif [ X"$1" = X"j2g" ]; then
		$GEN_COMMD $basepath/$JSON_SCRIPT $2 $3
	else
		echo "This type is not supported."
	fi
}

funcGen()
{
    if [ X"$1" = X"p2g" ]; then
        funcProto2go $2
        exit 0
    fi

	FULL_PATH=$curpath/$DIR_GEN
	if [ ! -d $FULL_PATH ]; then $MKDIR -p $FULL_PATH; fi;
	# 脚本类型全部转成小写
	typename=$(echo $1 | tr '[A-Z]' '[a-z]')

	if [ -d $2 ]; then
		for file in $2/*
		do
			file_name=${file##*/}
			if [ -f "$file" ]
			then
				if [ X"${file_name##*.}" = X"yml" -a X"$typename" = X"y2g" ] || [ X"${file_name##*.}" = X"json" -a X"$typename" = X"j2g" ]
				then
					des_file=${file_name%%.*}".go"
					gen $typename $file $FULL_PATH/$des_file
				fi
			fi
		done
	elif [ -f "$2" ]; then
		filename=$2
		file_name=${filename##*/}
		des_file=${file_name%%.*}".go"
		gen $1 $2 $FULL_PATH/${des_file##*/}
	fi
}

funcProtoRegister()
{
	fullpath="$(cd "$1" && pwd -P)"
	appname=`basename $fullpath`
	cat $1/*.go | awk -v appname=$appname '
	BEGIN{
		print "package "appname
		print
		print "import boot \"github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot\""
		print
		print "var _ boot.ServerType"
		print
		print "func init() {"
	}
	{
		if ($1=="type" && $3=="interface" && match($2, /.*Client/) && !match($2, /_/)) {
			print "	boot.RegisterClient(\""appname"\", New"$2")"
		}
	}
	END{
		print "}"
	}' > $fullpath/init.pb.go
}

funcProto2go()
{
  if [ "1$SKIP_P2G" = "1yes" ]; then
    echo "skip p2g"
    exit
  fi
	for dir in $(find $1 -type d); do
		if [ -z "$(ls $dir/*.proto 2>/dev/null)" ]; then
			# skip dir without .proto files
			continue
		fi
		protoc -I${ROOTPATH}/internal/common/proto --go_out=plugins=grpc:${ROOTPATH}/internal/common/proto --go_opt=paths=source_relative $dir/*.proto
		funcProtoRegister $dir
		for fl in `find $dir/*.pb.go`; do
		    protoc-go-inject-tag -XXX_skip=yaml,xml -input=$fl 2> /dev/null
		done
	done
}


if [ ! -n "$2" ]; then
	#funcHelp
	exit 1
fi

case $1 in
	create)
		funcCreate $2
		;;
	gen)
		if [ ! -n "$3" ]; then
			#funcHelp
			exit 1
		else
			funcGen $2 $3
		fi
		;;
	parser)
		funcParser $2
		;;
	help)
		funcHelp
		;;
	*)
		funcHelp
		exit 1
		;;
esac
