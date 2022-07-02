pkgname=promodj
pkgver=0.2.0
pkgrel=0
pkgdesc="Proxy client for promodj.com"
arch=('x86_64' 'aarch64')
license=('GPL')
url='https://github.com/nxshock/$pkgname'
depends=('ffmpeg')
makedepends=('go' 'git')
options=("!strip")
backup=("etc/$pkgname.conf")
source=("git+https://github.com/nxshock/$pkgname.git")
sha256sums=('SKIP')

build() {
	cd "$srcdir/$pkgname"
	go build -o $pkgname -buildmode=pie -trimpath -ldflags="-linkmode=external -s -w"
}

package() {
	cd "$srcdir/$pkgname"
	install -Dm755 "$pkgname"          "$pkgdir/usr/bin/$pkgname"
	install -Dm644 "$pkgname.conf"     "$pkgdir/etc/$pkgname.conf"
	install -Dm644 "$pkgname.service"  "$pkgdir/usr/lib/systemd/system/$pkgname.service"
	install -Dm644 "$pkgname.sysusers" "$pkgdir/usr/lib/sysusers.d/$pkgname.conf"
}
