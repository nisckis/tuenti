#!/usr/bin/env bash

function chr() {
    printf \\$(printf '%03o' $1)
}

function hex() {
    printf '%02X\n' $1
}

function encrypt() {
    key=$1
    msg=$2

    crpt_msg=""

    for ((i=0; i<${#msg}; i++)); do
        echo $i

        # The i-th char of the message
        c=${msg:$i:1}

        # Echo the character
        #  -n -> do not output the trailing newline
        #  -e -> enable interpretation of backslash escapes
        # pipe it to od
        # -A, --address-radix=RADIX
        #     output format for file offsets; RADIX is one of [doxn], 
        #     for Decimal, Octal, Hex or None
        # -t is printing format
        #     u[SIZE]
        #             unsigned decimal, SIZE bytes per integer
        #     c      printable character or backslash escape
        #
        # Basically the ASCII code of c
        asc_chr=$(echo -ne "$c" | od -An -tuC)

        # Compute a key position relative
        # to the msg len
        key_pos=$((${#key} - 1 - ${i}))

        # Get the char at the position
        # key_pos of the key
        key_char=${key:$key_pos:1}

        echo $c $asc_chr $key_pos $key_char
        echo "?" $(($asc_chr)) $((${key_char})) $(( $asc_chr ^ ${key_char} ))

        # Bitwise XOR
        crpt_chr=$(( $asc_chr ^ ${key_char} ))
        echo "XOR" $asc_chr $key_char $crpt_chr 

        # Call Hex function
        # hx_crpt_chr=$(hex $crpt_chr)
        hx_crpt_chr=$(hex $crpt_chr)
        echo "hex" $hx_crpt_chr
        echo "cipher ->" $crpt_msg "+" $hx_crpt_chr
        
        crpt_msg=${crpt_msg}${hx_crpt_chr}
        
        echo ""
    done
    echo $crpt_msg
}

# encrypt "12345" "message"
encrypt "40614178165780923111223" "514;248;980;347;145;332"
# decrypt "514;248;980;347;145;332" "3633363A33353B393038383C363236333635313A353336"
