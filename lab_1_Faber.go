package main

import "fmt"
import "os"
import "io"
import geojson "github.com/paulmach/go.geojson"
import "github.com/fogleman/gg"

//Структура данных для хранения параметров изображения 
type forProperties struct {
	Type         string
	Stroke       string
	Width        float64
	Opacity      float64
	Fill         string
	Fill_opacity float64
}

//Структура данных для хранения координат 
type position struct {
	X float64
	Y float64
}

//Переменная для хранения данных из файла geojson в байтовом представлении
var data64 = []byte(ReadFile("map.geojson"))
//Переменная для хранения для хранения параметров изображения
var Properties = forProperties{}
//Массив для хранения координат изображения
var Coordinates [255]position
//Переменная для хранения количества координат изображения
var NumCoordinates = 0

//Функция чтения данных из файла в строчном представлении
func ReadFile(geofile string) string {
	file, err := os.Open(geofile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data64 := make([]byte, 64)
	var data string

	for {
		fc, err := file.Read(data64)
		if err == io.EOF {
			break
		}
		data += string(data64[:fc])
	}

	return data
}

//Функция структурирования данных из файла geojson в байтовом представлении для последующего парсинга
func Collection() *geojson.FeatureCollection {
	fc, _ := geojson.UnmarshalFeatureCollection(data64)
	return fc
}

//Функция рендера изображения
func Drowing() {
	const width = 1366
	const height = 1024
	var coefW float64 = 0.003 * width
	var coefH float64 = 0.003 * height
	//Подготовка пространства вывода
 	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	//Отрисовка фигуры
	for coord := 0; coord < NumCoordinates-1; coord++ {
		x := (Coordinates[coord].X + 70) * coefW
		y := (Coordinates[coord].Y + 100) * coefH
		if coord == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.ClosePath()
	dc.SetHexColor(Properties.Fill)
	dc.Fill()

	//Отрисовка контура фигуры
	dc.SetHexColor(Properties.Stroke)
	dc.SetLineWidth(Properties.Width)
	for coord := 0; coord < NumCoordinates-1; coord++ {
		x := (Coordinates[coord].X + 70) * coefW
		y := (Coordinates[coord].Y + 100) * coefH
		if coord == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.ClosePath()
	dc.Stroke()

	//Сохранение в файл
	dc.SavePNG("out.png")
}

//Функция парсинга данных из файла geojson
func Parsing() {
	geojson_properties := Collection()

	//Парсинг параметров
	Properties.Stroke = geojson_properties.Features[0].Properties["stroke"].(string)
	Properties.Width = geojson_properties.Features[0].Properties["stroke-width"].(float64)
	Properties.Opacity = geojson_properties.Features[0].Properties["stroke-opacity"].(float64)
	Properties.Fill = geojson_properties.Features[0].Properties["fill"].(string)
	Properties.Fill_opacity = geojson_properties.Features[0].Properties["fill-opacity"].(float64)

	//Парсинг координат
	CoordinatesJSON := geojson_properties.Features[0].Geometry.Polygon
	NumCoordinates = len(CoordinatesJSON[0])
	for coord := 0; coord < NumCoordinates; coord++ {
		Coordinates[coord].X = CoordinatesJSON[0][coord][0]
		Coordinates[coord].Y = CoordinatesJSON[0][coord][1]
	}
}

//Главная функция программы, реализующая функцию парсинга geojson-файла и функцию рендера 
func main() {
	Parsing()
	Drowing()
}
