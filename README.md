# filter-auth
Logger for OpenSMTPD failed authentication attempts using OpenSMTPD reporting API.

The logs can be processed with fail2ban to block hosts with excessive failed login attempts.

### Log output:

    May  1 03:27:52 host smtpd[161]: auth-reporter: failed authentication from user=user@domain.org address=1.2.3.4 host=host.domain.com
