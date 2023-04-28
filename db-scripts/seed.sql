USE dermatologie24;

INSERT INTO users (firstname,lastname,email,password) VALUES
    (
        'Max',
        'Mustermann',
        'max.mustermann@gmail.com',
        'max.mustermann@gmail.com'
    ), (
        'Manuel',
        'Grabher',
        'manuel.grabher@web.de',
        'manuel.grabher@web.de'
    ),
    (
        'John',
        'Doe',
        'john.doe@outlook.com',
        'john.doe@outlook.com'
    )
    ;

INSERT INTO booking_status(status) VALUES("created"),("paid"),("completed"),("rejected");

INSERT INTO bookings (user_id, statusId, subject, message) VALUES
    (
        0,
        1,
        'Akne',
        'Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet'), (
        1,
    (
        1,
        1,
        'Akne',
        'Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet');


INSERT INTO booking_files (booking_id, file_path, name) VALUES (4, "https://www.oekotest.de/static_files/images/article/Mittel-gegen-Akne-im-Test_New-Africa-Shutterstock_106865_16x9.jpg", "akne.jpg"), (4, "https://image.brigitte.de/10915496/t/Yk/v4/w1440/r1.5/-/rachel-crawley-vorher-nachher.jpg", "akne2.jpg"), (5, "https://www.alterneudenken.com/wp-content/uploads/2022/09/neuroderm_beitrag2.jpg", "neurodermitis.jpg"), (5, "https://www.medipharma.de/Medipharma/MediLife/mediRatgeber/image-thumb__2474__WebTextImageMP/iStock-1091737400.jpg", "neurodermitis.jpg"), (6, "https://media.springernature.com/lw411/springer-cms/rest/v1/img/24020740/v7/4by3?as=jpg", "akne.jpg")