cmake_minimum_required(VERSION 3.9)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ../../dll)

project(
        Project.My.Hello
        VERSION 100
        DESCRIPTION "My.Hello"
        LANGUAGES C
)

add_library(My.Hello SHARED
        native.h native.c
        hello.c
)

target_include_directories(Nar.Hello PRIVATE ../narc)
set_target_properties(My.Hello PROPERTIES VERSION ${PROJECT_VERSION})

set(CMAKE_SHARED_LINKER_FLAGS "-I/opt/local/include")

target_link_libraries(My.Hello)
