# Отчет по дисциплине "Бэкенд разработка веб-приложений"

## 1. Название проекта
Flashcard Learning Backend API

## 2. Описание предметной области
Предметная область проекта - цифровая платформа для обучения по карточкам.
Пользователь регистрируется, авторизуется и управляет своими наборами карточек.
Внутри коллекции создаются карточки формата вопрос/ответ.
Также доступен режим тренировки, в котором карточки выдаются в случайном порядке.

Основная сущность для CRUD: коллекция.

## 3. Технологический стек
- Язык: Go
- HTTP роутер: Chi
- ORM: GORM
- База данных: PostgreSQL 17
- Авторизация: JWT
- Валидация: go-playground/validator
- Документация API: Swagger (swaggo, автогенерация)

## 4. Архитектура проекта
Используется слоистая архитектура через Dependency Injection.
Handler -> Service -> Repository -> Database

Кратко по слоям:
- Handler: парсинг запросов, валидация, HTTP-коды и JSON-ответы
- Service: бизнес-логика и проверки доступа
- Repository: работа с PostgreSQL через GORM
- DB: инициализация подключения, миграции, пул соединений

## 5. Реализованные API ручки
[Ручки](./swagger.yaml) можно посмотреть через онлайн [swagger-viewer](https://editor.swagger.io/): открыть сайт, затем File -> Import File и выбрать `swagger.yaml` из корня проекта.

### 5.1 Авторизация
- POST /auth/signup
- POST /auth/login
- POST /auth/logout

### 5.2 Профиль пользователя
- GET /profile/
- PUT /profile/username
- PUT /profile/password

### 5.3 Коллекции (основная сущность)
- GET /collections/
- GET /collections/{collectionID}/
- POST /collections/
- PUT /collections/{collectionID}/
- PATCH /collections/{collectionID}/
- DELETE /collections/{collectionID}/
- GET /collections/{collectionID}/train

### 5.4 Карточки
- GET /collections/{collectionID}/cards/
- POST /collections/{collectionID}/cards/
- PUT /cards/{cardID}/
- PATCH /cards/{cardID}/
- DELETE /cards/{cardID}/

## 6. Обработка ошибок
Используется единый формат:
```json
{
  "error": "текст ошибки"
}
```

Покрыты статусы:
- 400 Bad Request
- 401 Unauthorized
- 403 Forbidden
- 404 Not Found
- 405 Method Not Allowed
- 422 Unprocessable Entity
- 500 Internal Server Error

## 7. Валидация
Валидация реализована через пакет `go-playground/validator` (`validator/v10`).

Реализованы проверки:
- обязательные поля (required)
- корректность ID (gt=0)
- отклонение неизвестных полей JSON
- проверка формата входного JSON

## 8. Документация API
Swagger генерируется автоматически из аннотаций в коде.

Просмотр документации:
- открыть [swagger-viewer](https://editor.swagger.io/)
- File -> Import File -> выбрать [swagger.yaml](swagger.yaml)

Генерация через Makefile:
- `make swagger`
- выполняемая команда: `go run github.com/swaggo/swag/cmd/swag@v1.8.1 init -g cmd/main.go -o . -ot yaml --parseInternal`

Команда запуска тестов через Makefile:
- `make test`
- выполняемая команда: `go test ./...`

## 10. Список скриншотов для отчета
### 10.1 Предметная область
- PLACEHOLDER: Вставить скрин описания предметной области (раздел 2)

### 10.2 Скриншоты по каждой ручке
Для каждой ручки приложить 2 скрина:
1) участок кода с реализацией ручки
2) выполнение ручки в Swagger/Postman

POST /auth/signup
handler для регистрации 
```go
// SignUp registers a new user.
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.CreateUserReq true "signup payload"
// @Success 201 {object} schemas.CreateUserResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.CreateUserReq

	if err := util.DecodeJSONRequest(r, &userSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	createdUser, err := h.services.SignUp(&userSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdUser)
}
```
- Успешная регистрация
![body запроса](image.png)
![ответ сервера и возможные варианты](img/image-1.png)

POST /auth/login
handler login
```go
// SignIn authenticates user and returns JWT.
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.SignInReq true "signin payload"
// @Success 200 {object} schemas.SignInResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.SignInReq

	if err := util.DecodeJSONRequest(r, &userSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resp, err := h.services.SignIn(&userSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
```
- Скрин успешного логина. body запроса такой же, как при регистрации
![успешный login с ответами](img/image-2.png)

получим ошибку с использованием некорректного метода
![запрос и ответ](img/image-35.png)


POST /auth/logout
handler logout
```go
// LogOut is a formal logout endpoint for JWT architecture.
// @Summary Logout user
// @Tags auth
// @Produce json
// @Success 200 {object} schemas.LogOutResp
// @Router /auth/logout [post]
func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, schemas.LogOutResp{Message: "logged out"})
}
```
- Скрин успешного логаута
![логаут](img/image-3.png)


GET /profile/
handler getProfile
```go
// getProfile returns current user profile.
// @Summary Get profile
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} schemas.GetProfileResp
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /profile/ [get]
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	resp, err := h.services.GetProfile(userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
```
- Скрин запроса/ответа
![запрос информации о профиле](img/image-4.png)

PUT /profile/username
handler updateUsername
```go
// updateUsername updates current user username.
// @Summary Update username
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.UpdateUsernameReq true "username payload"
// @Success 200 {object} schemas.UpdateUsernameResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /profile/username [put]
func (h *Handler) updateUsername(w http.ResponseWriter, r *http.Request) {
	var updateUsernameReq schemas.UpdateUsernameReq

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	updateUsernameReq.ID = userID

	if err := util.DecodeJSONRequest(r, &updateUsernameReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&updateUsernameReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resp, err := h.services.UpdateUsername(&updateUsernameReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
```
- Скрин запроса/ответа
![запрос](img/image-5.png)
![ответ](img/image-6.png)

PUT /profile/password
handler updatePassword
```go
// updatePassword updates current user password.
// @Summary Update password
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.UpdatePasswordBody true "password payload"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /profile/password [put]
func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	var updatePasswordBody schemas.UpdatePasswordBody
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := util.DecodeJSONRequest(r, &updatePasswordBody); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	updatePasswordReq := schemas.UpdatePasswordReq{
		ID:          userID,
		OldPassword: updatePasswordBody.OldPassword,
		NewPassword: updatePasswordBody.NewPassword,
	}

	if err := h.validator.ValidateWithDetailedErrors(&updatePasswordReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.services.UpdatePassword(&updatePasswordReq); err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusNoContent, nil)
}
```
- Скрин запроса/ответа
![запрос](img/image-7.png)
![ответ](img/image-8.png)


GET /collections/
handler getAllCollections
```go
// getAllCollections returns all current user collections.
// @Summary Get all collections
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} schemas.AllCollectionsResp
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /collections/ [get]
func (h *Handler) getAllCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	allCollections, err := h.services.GetAllCollections(userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, allCollections)
}
```
- PLACEHOLDER: Скрин запроса/ответа
![запрос и ответ](img/image-9.png)

GET /collections/{collectionID}/
handler getCollectionByID
```go
// getCollectionByID returns collection by ID.
// @Summary Get collection by ID
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.GetCollectionByIDResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [get]
func (h *Handler) getCollectionByID(w http.ResponseWriter, r *http.Request) {
	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	collection, err := h.services.GetCollectionByID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, collection)
}
```
- Скрин запроса/ответа
![запрос и ответ](img/image-10.png)


POST /collections/
handler createCollection
```go
// createCollection creates a new collection for the current user.
// @Summary Create collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.CreateCollectionReq true "collection payload"
// @Success 201 {object} schemas.CreateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /collections/ [post]
func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionSchemaReq schemas.CreateCollectionReq

	if err := util.DecodeJSONRequest(r, &collectionSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&collectionSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	createdCollection, err := h.services.CreateCollection(&collectionSchemaReq, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdCollection)
}
```
- PLACEHOLDER: Скрин успешного создания коллекции
![запрос](img/image-11.png)
![ответ](img/image-12.png)

PUT /collections/{collectionID}/
handler editCollection
```go
// editCollection fully updates collection.
// @Summary Full update collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.UpdateCollectionReq true "collection payload"
// @Success 200 {object} schemas.UpdateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [put]
func (h *Handler) editCollection(w http.ResponseWriter, r *http.Request) {
	var updatedCollectionSchema schemas.UpdateCollectionReq

	if err := util.DecodeJSONRequest(r, &updatedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	updatedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCollection, err := h.services.UpdateCollection(&updatedCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, updatedCollection)
}
```
- PLACEHOLDER: Скрин полного обновления коллекции
![запрос](img/image-15.png)
![ответ](img/image-14.png)

PATCH /collections/{collectionID}/
handler patchCollection
```go
// patchCollection partially updates collection.
// @Summary Partial update collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.PatchCollectionReq true "collection payload"
// @Success 200 {object} schemas.UpdateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [patch]
func (h *Handler) patchCollection(w http.ResponseWriter, r *http.Request) {
	var patchCollectionSchema schemas.PatchCollectionReq

	if err := util.DecodeJSONRequest(r, &patchCollectionSchema); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	patchCollectionSchema.ID = collectionID

	if patchCollectionSchema.Name == nil && patchCollectionSchema.Description == nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidRequestBody)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&patchCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCollection, err := h.services.PatchCollection(&patchCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, updatedCollection)
}
```
- PLACEHOLDER: Скрин частичного обновления коллекции
![запрос](img/image-16.png)
![ответ](img/image-17.png)

DELETE /collections/{collectionID}/
handler removeCollection
```go
// removeCollection deletes collection.
// @Summary Delete collection
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [delete]
func (h *Handler) removeCollection(w http.ResponseWriter, r *http.Request) {
	var removedCollectionSchema schemas.RemoveCollectionReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	removedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&removedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.services.RemoveCollection(&removedCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusNoContent, nil)
}
```
- PLACEHOLDER: Скрин удаления коллекции
![запрос](img/image-18.png)
![ответ](img/image-19.png)

GET /collections/{collectionID}/train
handler startPractise
```go
// startPractise returns shuffled cards for training.
// @Summary Train cards
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.TrainSchemaResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/train [get]
func (h *Handler) startPractise(w http.ResponseWriter, r *http.Request) {
	var practiseSchemaReq schemas.TrainSchemaReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	practiseSchemaReq.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&practiseSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	randomCards, err := h.services.TrainCards(&practiseSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, randomCards)
}
```
- PLACEHOLDER: Скрин режима тренировки
![вопрос и ответ](img/image-20.png)

GET /collections/{collectionID}/cards/
handler getCardsByCollectionID
```go
// getCardsByCollectionID returns all cards from collection.
// @Summary Get cards by collection
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.GetCardByCollectionIDResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /collections/{collectionID}/cards/ [get]
func (h *Handler) getCardsByCollectionID(w http.ResponseWriter, r *http.Request) {
	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	cards, err := h.services.GetCardsByCollectionID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, cards)
}
```
- PLACEHOLDER: Скрин запроса/ответа
![вопрос и ответ](img/image-21.png)

чтобы получить 404, запросим несуществующиее значение
![запрос и ответ](img/image-34.png)

POST /collections/{collectionID}/cards/
handler createCard
```go
// createCard creates a card in collection.
// @Summary Create card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.CreateCardReq true "card payload"
// @Success 201 {object} schemas.CreateCardResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/cards/ [post]
func (h *Handler) createCard(w http.ResponseWriter, r *http.Request) {
	var cardSchemaReq schemas.CreateCardReq

	if err := util.DecodeJSONRequest(r, &cardSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&cardSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	createdCard, err := h.services.CreateCard(&cardSchemaReq, collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdCard)
}
```
- PLACEHOLDER: Скрин успешного создания карточки
![запрос](img/image-22.png)
![ответ](img/image-23.png)

PUT /cards/{cardID}/
handler editCard (для PUT)
```go
// editCard updates card (used for PUT and PATCH routes).
// @Summary Update card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param cardID path int true "card id"
// @Param request body schemas.UpdateCardReq true "card payload"
// @Success 200 {object} schemas.UpdateCardResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /cards/{cardID}/ [put]
// @Router /cards/{cardID}/ [patch]
func (h *Handler) editCard(w http.ResponseWriter, r *http.Request) {
	var updatedCardSchemaReq schemas.UpdateCardReq

	if err := util.DecodeJSONRequest(r, &updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCardAccess(cardID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	updatedCardSchemaReq.ID = cardID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCard, err := h.services.UpdateCard(&updatedCardSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, updatedCard)
}
```
- PLACEHOLDER: Скрин полного обновления карточки
![запрос](img/image-24.png)
![ответ](img/image-25.png)

PATCH /cards/{cardID}/
handler editCard (для PATCH)
```go
// editCard updates card (used for PUT and PATCH routes).
// @Summary Update card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param cardID path int true "card id"
// @Param request body schemas.UpdateCardReq true "card payload"
// @Success 200 {object} schemas.UpdateCardResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /cards/{cardID}/ [put]
// @Router /cards/{cardID}/ [patch]
func (h *Handler) editCard(w http.ResponseWriter, r *http.Request) {
	var updatedCardSchemaReq schemas.UpdateCardReq

	if err := util.DecodeJSONRequest(r, &updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCardAccess(cardID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	updatedCardSchemaReq.ID = cardID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCard, err := h.services.UpdateCard(&updatedCardSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, updatedCard)
}
```
- PLACEHOLDER: Скрин частичного обновления карточки
![запрос](img/image-26.png)
![ответ](img/image-27.png)

DELETE /cards/{cardID}/
handler removeCard
```go
// removeCard deletes card.
// @Summary Delete card
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param cardID path int true "card id"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /cards/{cardID}/ [delete]
func (h *Handler) removeCard(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCardAccess(cardID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	removeCardReq := schemas.RemoveCardReq{ID: cardID}

	if err := h.validator.ValidateWithDetailedErrors(&removeCardReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.services.RemoveCard(&removeCardReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusNoContent, nil)
}
```
- PLACEHOLDER: Скрин удаления карточки
![вопрос и ответ](img/image-28.png)

### 10.3 Неавторизованный доступ
middleware-проверка токена
```go
func UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrUnauthorized)
			return
		}

		claims, err := util.ParseToken(headerParts[1])
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
			return
		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```
- PLACEHOLDER: Скрин запроса к защищенной ручке без токена
![запрос и ответ](img/image-29.png)


### 10.4 Валидация
код validator/v10
```go
func (v *Validator) ValidateWithDetailedErrors(i interface{}) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, validationError := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' %s", validationError.Field(), validationError.ActualTag()))
		}
		return fmt.Errorf("invalid input: %s", strings.Join(errorMessages, ", "))
	}

	return err
}
```
Нельзя отправить некорректный тип в swagger, поэтому запрос отправлен через терминал
![запрос](img/image-30.png)

строгий JSON-декодер
```go
func DecodeJSON(schema interface{}, raw []byte) error {
	if len(raw) == 0 {
		return io.ErrUnexpectedEOF
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(schema); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return &json.SyntaxError{Offset: 0}
	}

	return nil
}
```
- PLACEHOLDER: Скрин ответа 422 на невалидные данные
Опустим обязательный ключ "name", получим ошибку 422
![запрос](img/image-31.png)
![ответ](img/image-32.png)
- PLACEHOLDER: Скрин ответа 400 на некорректный JSON
Сломаем json с помощью лишней запятой в конце
![запрос и ответ](img/image-33.png)

### 10.6 Дополнительные материалы
- PLACEHOLDER: Приложить README
- PLACEHOLDER: Приложить env/config без секретов
- PLACEHOLDER: Скрин подтверждения данных в БД преподавателя

## 11. Заключение
В ходе работы реализован REST API с авторизацией, CRUD, обработкой ошибок и валидацией по требованиям задания.
Проект документирован (README + Swagger), покрыт тестами и готов к демонстрации преподавателю.
