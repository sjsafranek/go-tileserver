
import sys
from osgeo import ogr
import os


shapefile = sys.argv[1]

datasource = ogr.Open(shapefile)
# driver = ogr.GetDriverByName("ESRI Shapefile")
# dataSource = driver.Open(shapefile, 0)
layer = datasource.GetLayer(0)
# layer = dataSource.GetLayer()
layerDefinition = layer.GetLayerDefn()


string_fields = []

print("Name\t\tType\tWidth\tPrecision")
for i in range(layerDefinition.GetFieldCount()):
    fieldName =  layerDefinition.GetFieldDefn(i).GetName()
    fieldTypeCode = layerDefinition.GetFieldDefn(i).GetType()
    fieldType = layerDefinition.GetFieldDefn(i).GetFieldTypeName(fieldTypeCode)
    fieldWidth = layerDefinition.GetFieldDefn(i).GetWidth()
    GetPrecision = layerDefinition.GetFieldDefn(i).GetPrecision()

    print( fieldName + "\t\t" + fieldType+ "\t" + str(fieldWidth) + "\t" + str(GetPrecision) )

    if 'String' == fieldType:
        string_fields.append(fieldName)


# for feature in layer:
#     for field in string_fields:
#         print(feature.GetField(field))
# layer.ResetReading()

datasource.Destroy()

