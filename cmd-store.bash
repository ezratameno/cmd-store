#! /bin/bash


_completions() {
  local cur prev
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=( $(compgen -W "aws s3 " -- "$cur") )

  
  elif [[ ${COMP_WORDS[1]} == "$aws" ]]; then
    COMPREPLY=( $(compgen -W " pc list-buckets file" -- "$cur") )
  
  elif [[ ${COMP_WORDS[1]} == "$s3" ]]; then
    COMPREPLY=( $(compgen -W " list-buckets" -- "$cur") )
  
  fi
}


complete -F _completions ./cmd-store 