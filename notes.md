se debe crear un module para poder instalar paquetes
`go mod init github.com/masmerino13/lens/v2`

`go.mod` contiene los paquetes instalados

limpia o reset do.mod, borra las versiones de los paquetes agregados al module
`go mod tidy`


uso de `...` variadic parameters
se utiliza para pasar argumentos de forma dinamica a las funciones

```go
func main() {
  Demo(1, 2, 3) // se pasan a la funcion 3 numeros

  nums := []int{1, 2, 3}

  Demo(nums) // Falla debido a que Demo espera varios argumentos

  Demo(nums...) // funciona, los 3 puntos a la derecha resuelven el problema
}

func Demo(numbers ...int) {
  // do
}

```

