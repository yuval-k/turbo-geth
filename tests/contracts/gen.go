package contracts

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build constructor_and_var.sol
//go:generate abigen -abi build/constructor_and_var.abi -bin build/constructor_and_var.bin -pkg contracts -type constructor_and_var -out ./gen_constructor_and_var.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build constructor_and_var_usage.sol
//go:generate abigen -abi build/constructor_and_var_usage.abi -bin build/constructor_and_var_usage.bin -pkg contracts -type constructor_and_var_usage -out ./gen_constructor_and_var_usage.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build only_constructor.sol
//go:generate abigen -abi build/only_constructor.abi -bin build/only_constructor.bin -pkg contracts -type only_constructor -out ./gen_only_constructor.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build simple_function.sol
//go:generate abigen -abi build/simple_function.abi -bin build/simple_function.bin -pkg contracts -type simple_function -out ./gen_simple_function.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build two_function.sol
//go:generate abigen -abi build/two_function.abi -bin build/two_function.bin -pkg contracts -type two_function -out ./gen_two_function.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build two_function_call_with_return.sol
//go:generate abigen -abi build/two_function_call_with_return.abi -bin build/two_function_call_with_return.bin -pkg contracts -type two_function_call_with_return -out ./gen_two_function_call_with_return.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build two_function_call_without_return.sol
//go:generate abigen -abi build/two_function_call_without_return.abi -bin build/two_function_call_without_return.bin -pkg contracts -type two_function_call_without_return -out ./gen_two_function_call_without_return.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build two_function_empty.sol
//go:generate abigen -abi build/two_function_empty.abi -bin build/two_function_empty.bin -pkg contracts -type two_function_empty -out ./gen_two_function_empty.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build with_if.sol
//go:generate abigen -abi build/with_if.abi -bin build/with_if.bin -pkg contracts -type with_if -out ./gen_with_if.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build with_if_continue.sol
//go:generate abigen -abi build/with_if_continue.abi -bin build/with_if_continue.bin -pkg contracts -type with_if_continue -out ./gen_with_if_continue.go

//go:generate solc --allow-paths ., --abi --bin --overwrite -o build with_if_else.sol
//go:generate abigen -abi build/with_if_else.abi -bin build/with_if_else.bin -pkg contracts -type with_if_else -out ./gen_with_if_else.go
