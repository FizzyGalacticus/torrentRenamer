# torrentRenamer

## DISCLAIMER

This project in no way endorses copyright infringement, and should only be used on videos that you legally own or otherwise have the right to have.

## About

`torrentRenamer` uses the [parse-torrent-name](https://github.com/middelink/go-parse-torrent-name) project to parse video filenames and rename them for better organization. With the configuration options, it is very easy to set default paths for movies and shows.

By default, videos are renamed to:

-   Movies: `/path/to/home/Videos/Movies/{title} - ({year}).{extension}`
-   TV Shows: `/path/to/home/TV Shows/{series}/{series} - Season {season}/{series} - S{season}E{episode}.{extension}`

## OMDB Support

There is support to use the [OMDB API](http://www.omdbapi.com/). To use it, all you will need to do is set the `--omdb-key` flag.

When OMDB is used, it will use the `title`, `season`, `episode`, and `year` values that are returned from it, instead of from `parse-torrent-name`. It will also change the output filename of a TV Show to `{series} = S{season}E{episode} - {title}.{extension}` (notice the added episode title).

[Get your **free** API key here](https://www.omdbapi.com/apikey.aspx)

## Options

| Option Name           | Usages                  | Defaults                                                                                                                                                          |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Movies Directory      | `--movies`\|`-m`        | <home_dir>/Videos/Movies                                                                                                                                          |
| Shows Directory       | `--shows`\|`-s`         | <home_dir>/Videos/TV Shows                                                                                                                                        |
| Movie Template        | `--movie-template|`     | `"{{ .Name }} ({{ .Year }}).{{ .Ext }}"`                                                                                                                          |
| Show Template         | `--show-template`       | `"{{ .Name }}{{sep}}{{ .Name }} - Season {{padDigit .Season 2}}{{sep}}{{ .Name }} - S{{padDigit .Season 2 }}E{{padDigit .Episode 2 }}.{{ .Ext }}"`                |
| Service               | `--service`             | `nil`                                                                                                                                                             |
| OMDB API Key          | `--omdb-key`            | `nil`                                                                                                                                                             |
| OMDB Movie Template   | `--omdb-movie-template` | `"{{ .Name }} ({{ .Year }}).{{ .Ext }}"`                                                                                                                          |
| OMDB Show Template    | `--omdb-show-template`  | `"{{ .Name }}{{sep}}{{ .Name }} - Season {{padDigit .Season 2}}{{sep}}{{ .Name }} - S{{padDigit .Season 2 }}E{{padDigit .Episode 2 }} - {{ .Title }}.{{ .Ext }}"` |
| Add Name Override     | `--add-override`        | `nil`                                                                                                                                                             |
| Remove Name Override  | `--rm-override`         | `nil`                                                                                                                                                             |
| Save Config           | `--save-config`         | `false`                                                                                                                                                           |
| Rename Without Prompt | `--yes`                 | `-y`                                                                                                                                                              | `false` |

### Notes

#### Templates

The template follows the same format as [Go's text/template package](https://golang.org/pkg/text/template/). There are currently a few custom functions available:

* `padDigit` - Takes two integers, and makes sure that the first one passed contains at least as many digits as specified with the second. If the first has too few, it will prepend zeroes.
  * Example: `"{{padDigit 2 2}}"` will result in `"02"`
* `sep` - Returns the OS specific path separator.
* `home` - Returns the users home directory on the running platform.
* `homePath` - Takes a path, and prepends the users home directory to it.

Apart from those functions, there are a handful of variables you can use for videos:

* `.Name` - The name of the movie or show
* `.Year` - The year of the movie.
  * **only works with movies**
* `.Season` - The season number
  * **only works with shows**
* `.Episode` - The episode number
  * **only works with shows**
* `.Title` - The title of the episode
  * **only works with shows and OMDB integration**
* `.Ext` - The file extension

You can permanently save/change your default template for the given category by setting it when running with the `--save-config` flag, or by manually modifying the `<home_dir>/.torrentRenamerrc` file.
