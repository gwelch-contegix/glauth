%global debug_package %{nil}

Name:           glauth
Version:        3.4.8
Release:        83%{?dist}
Summary:        glauth

License:        MIT
URL:            https://github.com/gwelch-contegix/glauth
Source0:        glauth.tar.gz
Source1:        glauth.service
Source2:        ldap-server

BuildRequires:  go-toolset
BuildRequires:  systemd-rpm-macros
BuildRequires:  make

%description
glauth

%prep
%setup -n v2


%build
make M=pkg/plugins/glauth-keycloak pull-plugin-dependencies
make releasemain
make P=keycloak M=pkg/plugins/glauth-keycloak releaseplugin
sed 's@/Users/gwelch/build/source/glauth/v2/bin/darwinamd64@%{_libdir}@g' pkg/plugins/glauth-keycloak/sample-keycloak.cfg >bin/glauth.cfg


%install
rm -rf $RPM_BUILD_ROOT
install -D -m 755 bin/%{name} %{buildroot}/%{_bindir}/%{name}
install -D -m 755 bin/keycloak.so %{buildroot}/%{_libdir}/keycloak.so
install -D -m 600 bin/glauth.cfg %{buildroot}/%{_sysconfdir}/glauth.cfg
install -D -m 644 %{SOURCE1} %{buildroot}/%{_unitdir}/%{name}.service


%files
%{_bindir}/%{name}
%{_libdir}/keycloak.so
%{_unitdir}/%{name}.service
%config %{_sysconfdir}/glauth.cfg
%license LICENSE


#%files cron
#%config %{_sysconfdir}/cron.weekly/%{name}.sh


%changelog
* Mon Aug  7 2023 gwelch
-
