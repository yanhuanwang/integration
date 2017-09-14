#Robot functions that will be shared also with other tests

*** Keywords ***
json_from_file
#Robot function to extract the json object from a file
    [Arguments]    ${file_path}
    ${json_file}=    Get file    ${file_path}
    ${json_object}=    Evaluate    json.loads('''${json_file}''')    json
    [return]    ${json_object}

string_from_json
#Robot function to transform the json object to a string
    [Arguments]    ${json_value}
    ${json_string}=   Stringify Json     ${json_value}
    [return]    ${json_string}

random_ip
#Robot function to generate a random IP
    [Arguments]
    ${numbers}=    Evaluate    random.sample([x for x in range(1, 256)], 4)    random
    ${generated_ip}=    Catenate    ${numbers[0]}.${numbers[1]}.${numbers[2]}.${numbers[3]}
    [return]    ${generated_ip}