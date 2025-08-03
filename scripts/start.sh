if [ ! -f ../cmd/cli/cli ]
then
    echo "Service not found"
    cd ../cmd/cli
    go build
    cd ../../scripts
fi
../cmd/cli/cli -file_in=../data/phrases/phrases1.txt -file_out=phrases1_out.txt
#natasha_tst.exe -file_in=phrases1_1.txt -file_out=phrases1_1_out.txt
#natasha_tst.exe -file_in=phrases1_2.txt -file_out=phrases1_2_out.txt
#natasha_tst.exe -file_in=phrases1_3.txt -file_out=phrases1_3_out.txt
#natasha_tst.exe -file_in=phrases1_4.txt -file_out=phrases1_4_out.txt
#natasha_tst.exe -file_in=phrases1_5.txt -file_out=phrases1_5_out.txt
#natasha_tst.exe -file_in=phrases1_6.txt -file_out=phrases1_6_out.txt
#natasha_tst.exe -file_in=phrases1_7.txt -file_out=phrases1_7_out.txt
#natasha_tst.exe -file_in=phrases1_8_1.txt -file_out=phrases1_8_1_out.txt
#natasha_tst.exe -file_in=phrases1_8_2.txt -file_out=phrases1_8_2_out.txt
#natasha_tst.exe -file_in=phrases1_8_3.txt -file_out=phrases1_8_3_out.txt
#natasha_tst.exe -file_in=phrases1_8_4.txt -file_out=phrases1_8_4_out.txt
#natasha_tst.exe -file_in=phrases1_8_5.txt -file_out=phrases1_8_5_out.txt
#natasha_tst.exe -file_in=phrases1_8_6.txt -file_out=phrases1_8_6_out.txt
#natasha_tst.exe -file_in=phrases1_8_7.txt -file_out=phrases1_8_7_out.txt
#natasha_tst.exe -file_in=phrases1_8_8.txt -file_out=phrases1_8_8_out.txt
#natasha_tst.exe -file_in=phrases1_short.txt -file_out=phrases1_short_out.txt
