(function($) {
	$.fn.catchEnter = function(sel) {
		return this.each(function() {
			$(this).on('keyup', sel, function(e) {
				if (e.ctrlKey && e.keyCode == 13)
					$(this).trigger("enterkey");
			})
		});
	};
})(jQuery);

$(document).ready(function() {

	// attach event
	$("#submit-btn").on('click', addSubmission);
	$("#content-box").catchEnter().on('enterkey', addSubmission);

	$("#create-btn").on('click', createSubmission);
	$("#follow-btn").on('click', followSubmissions);
	$("#content-box").on('keydown', function() {
		$("#resp").text('Write to...')
	});

});



async function createSubmission() {
	const data = {
        chatDesc: $("#chat-name").val(),
    };

    await fetch('/api/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then((response) => response.json())
        .then((resp) => {
            window.location.href = `/chats/${resp.chatDesc}`;
        })
        .catch((error) => {
            $("#resp").text(error);
        });
}

async function followSubmissions() {
	const data = {
        chatDesc: window.location.pathname.split("/").pop(),
    };

    await fetch('/api/follow', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .catch((error) => {
            $("#resp").text(error);
        });
}

async function addSubmission() {

	$("#submit-btn").prop("disabled", true); // disable multiple submission

	// prepare alert
	let card = $("#resp");

	// validate
	let content = $("#content-box").val();

	if ($.trim(content) === '') {
		$("#submit-btn").prop("disabled", false);
		card.text("Please type in your story first!");

		return;
	}

	const data = {
		chatDesc: window.location.pathname.split("/").pop(),
		message: content
	};

	await fetch('/api/submit', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		})
		.then((response) => response.json())
		.then((resp) => {
			card.text(resp.message);
		})
		.catch((error) => {
			card.text(error);
		});

	$("#submit-btn").prop("disabled", false);
}