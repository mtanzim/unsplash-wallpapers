# unsplash-wallpapers

A CLI tool that helps download images from Unsplash collections.

## Usage

```bash
go build
./unsplash-wallpapers -h
# as an example:
./unsplash-wallpapers -c 44204348 -d "./images"
```

![Help](./assets/readme.PNG)

## Note

Please provide a `.env` file with the following to gain access to the Unsplash API. See the [documentation](https://unsplash.com/developers) for details.

```text
SECRET=
ACCESS=
```

A [fixture](./collections/collections.json) with the expected API response from Unsplash is provided for convenience. [quicktype](https://app.quicktype.io/) was used to auto generate the struct based on the example fixture.
